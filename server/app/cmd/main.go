package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	userv1 "github.com/BrianLusina/skillq/server/app/api/rest/routes/users/v1"
	"github.com/BrianLusina/skillq/server/app/cmd/config"
	"github.com/BrianLusina/skillq/server/app/internal/app"
	"github.com/BrianLusina/skillq/server/infra/clients/email"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqpconsumer "github.com/BrianLusina/skillq/server/infra/messaging/amqp/consumer"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	appLogger := logger.New()

	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		appLogger.Error("failed set max procs", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		appLogger.Error("failed get config", err)
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// TODO: setup config
	app := fiber.New(fiber.Config{
		ServerHeader: "SkillQ",
		AppName:      "SkillQ",
	})

	go func() {
		defer app.Server().Shutdown()
		<-ctx.Done()
	}()

	// middleware
	app.Use(cors.New())

	// prepare and setup app
	setupApp(ctx, cancel, app, cfg, appLogger)

	// Start the server
	err = app.Listen(fmt.Sprintf(":%d", cfg.HTTP.Port))
	appLogger.Fatalf("Failed to start application: %v", err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		appLogger.Info("signal.Notify", v)
	case done := <-ctx.Done():
		appLogger.Info("ctx.Done", "app done", done)
	}
}

func setupApp(ctx context.Context, cancel context.CancelFunc, app *fiber.App, cfg *config.Config, appLogger logger.Logger) {
	// configuration
	mongoDbConfig := mongodb.MongoDBConfig{
		Client: mongodb.ClientOptions{
			Host:        cfg.MongoDB.Host,
			Port:        cfg.MongoDB.Port,
			User:        cfg.MongoDB.User,
			Password:    cfg.MongoDB.Password,
			RetryWrites: cfg.MongoDB.RetryWrites,
		},
		DBConfig: mongodb.DatabaseConfig{
			DatabaseName: cfg.MongoDB.Database,
		},
	}

	userMongodbConfig := mongodb.MongoDBConfig{
		Client: mongoDbConfig.Client,
		DBConfig: mongodb.DatabaseConfig{
			DatabaseName:   mongoDbConfig.DBConfig.DatabaseName,
			CollectionName: "users",
		},
	}

	// userVerificationMongodbConfig := mongodb.MongoDBConfig{
	// 	Client: mongoDbConfig.Client,
	// 	DBConfig: mongodb.DatabaseConfig{
	// 		DatabaseName:   mongoDbConfig.DBConfig.DatabaseName,
	// 		CollectionName: "userVerification",
	// 	},
	// }

	amqpConfig := amqp.Config{
		Username: cfg.RabbitMQ.Username,
		Password: cfg.RabbitMQ.Password,
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
	}

	minioConfig := minio.Config{
		Endpoint:        cfg.MinioConfig.Endpoint,
		AccessKeyID:     cfg.MinioConfig.AccessKeyID,
		SecretAccessKey: cfg.MinioConfig.SecretAccessKey,
		UseSSL:          cfg.MinioConfig.UseSSL,
		Token:           cfg.MinioConfig.Token,
	}

	emailConfig := email.EmailClientConfig{
		Host:     cfg.EmailConfig.Host,
		Port:     cfg.EmailConfig.Port,
		Password: cfg.EmailConfig.Password,
		From:     cfg.EmailConfig.From,
	}

	skillQApp := prepareApp(ctx, cancel, userMongodbConfig, amqpConfig, minioConfig, emailConfig)

	// routing
	userApi := userv1.NewUserApi(skillQApp.UserSvc, appLogger)
	userApi.RegisterHandlers(app)
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, mongoDbConfig mongodb.MongoDBConfig, amqpConfig amqp.Config, minioConfig minio.Config, emailConfig email.EmailClientConfig) *app.App {
	app, err := app.InitApp(mongoDbConfig, amqpConfig, minioConfig, emailConfig)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
		<-ctx.Done()
	}

	app.StoreImageEventPublisher.Configure(
		amqppublisher.Exchange(
			amqp.ExchangeOptionParams{
				Name:    "store-image-exchange",
				Kind:    "fanout",
				Durable: true,
			},
		),
		amqppublisher.BindingKey("store-image-routing-key"),
	)

	// Configure publisher and start workers
	app.SendEmailEventPublisher.Configure(
		amqppublisher.Exchange(
			amqp.ExchangeOptionParams{
				Name:    "send-email-exchange",
				Kind:    "fanout",
				Durable: true,
			},
		),
		amqppublisher.BindingKey("send-email-routing-key"),
	)

	app.AmqpEventConsumer.Configure(
		amqpconsumer.Exchange(
			amqp.ExchangeOptionParams{
				Name:    "store-image-exchange",
				Kind:    "fanout",
				Durable: true,
			},
		),
		amqpconsumer.Queue(
			amqp.QueueOptionParams{
				Name: "store-image-queue",
			},
		),
		amqpconsumer.BindingKey("store-image-routing-key"),
		amqpconsumer.Consumer(
			amqp.ConsumerOptionParams{
				Tag: "store-image-consumer",
			},
		),
	)

	app.AmqpEventConsumer.Configure(
		amqpconsumer.Exchange(
			amqp.ExchangeOptionParams{
				Name:    "send-email-exchange",
				Kind:    "fanout",
				Durable: true,
			},
		),
		amqpconsumer.Queue(
			amqp.QueueOptionParams{
				Name: "send-email-queue",
			},
		),
		amqpconsumer.BindingKey("send-email-routing-key"),
		amqpconsumer.Consumer(
			amqp.ConsumerOptionParams{
				Tag: "send-email-consumer",
			},
		),
	)

	go func() {
		err1 := app.AmqpEventConsumer.StartConsumer(app.Worker)
		if err1 != nil {
			slog.Error("Failed to start app Consumer", err1)
			cancel()
			<-ctx.Done()
		}
	}()

	return app
}
