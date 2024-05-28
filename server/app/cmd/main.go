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
	userapp "github.com/BrianLusina/skillq/server/app/internal/app/user"
	verificationapp "github.com/BrianLusina/skillq/server/app/internal/app/verification"
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
	prepareApp(ctx, cancel, app, cfg, appLogger)

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

func prepareApp(ctx context.Context, cancel context.CancelFunc, app *fiber.App, cfg *config.Config, appLogger logger.Logger) {
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

	userVerificationMongodbConfig := mongodb.MongoDBConfig{
		Client: mongoDbConfig.Client,
		DBConfig: mongodb.DatabaseConfig{
			DatabaseName:   mongoDbConfig.DBConfig.DatabaseName,
			CollectionName: "userVerification",
		},
	}

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

	userApp := prepareUserApp(ctx, cancel, userMongodbConfig, amqpConfig, minioConfig)
	_ = prepareUserVerificationApp(ctx, cancel, userVerificationMongodbConfig, amqpConfig, minioConfig)

	// routing
	userApi := userv1.NewUserApi(userApp.UserSvc, appLogger)
	userApi.RegisterHandlers(app)
}

func prepareUserApp(ctx context.Context, cancel context.CancelFunc, mongoDbConfig mongodb.MongoDBConfig, amqpConfig amqp.Config, minioConfig minio.Config) *userapp.UserApp {
	userApp, err := userapp.InitApp(mongoDbConfig, amqpConfig, minioConfig)
	if err != nil {
		slog.Error("failed init user app", err)
		cancel()
		<-ctx.Done()
	}

	// Configure publisher and start workers
	userApp.AmqpEventPublisher.Configure(
		amqppublisher.ExchangeName("email-verification-exchange"),
		amqppublisher.BindingKey("email-verification-routing-key"),
		amqppublisher.MessageTypeName("email-verification-started"),
	)

	userApp.AmqpEventConsumer.Configure(
		amqpconsumer.ExchangeName("skillq-user-exchange"),
		amqpconsumer.QueueName("skillq-user-queue"),
		amqpconsumer.BindingKey("skillq-user-routing-key"),
		amqpconsumer.ConsumerTag("skillq-user-consumer"),
	)

	go func() {
		err1 := userApp.AmqpEventConsumer.StartConsumer(userApp.Worker)
		if err1 != nil {
			slog.Error("failed to start user app Consumer", err1)
			cancel()
			<-ctx.Done()
		}
	}()

	return userApp
}

func prepareUserVerificationApp(ctx context.Context, cancel context.CancelFunc, mongoDbConfig mongodb.MongoDBConfig, amqpConfig amqp.Config, minioConfig minio.Config) *verificationapp.UserVerificationApp {
	userVerificationApp, err := verificationapp.InitializeUserVerificationApp(mongoDbConfig, amqpConfig, minioConfig)
	if err != nil {
		slog.Error("failed init verification app", err)
		cancel()
		<-ctx.Done()
	}

	// Configure publisher and start workers
	userVerificationApp.AmqpEventPublisher.Configure(
		amqppublisher.ExchangeName("skillq-user-verification-exchange"),
		amqppublisher.BindingKey("skillq-user-verification-routing-key"),
		amqppublisher.MessageTypeName("skillq-user-verification"),
	)

	userVerificationApp.AmqpEventConsumer.Configure(
		amqpconsumer.ExchangeName("skillq-user-verification-exchange"),
		amqpconsumer.QueueName("skillq-user-verification-queue"),
		amqpconsumer.BindingKey("skillq-user-verification-key"),
		amqpconsumer.ConsumerTag("skillq-user-verification-consumer"),
	)

	go func() {
		err1 := userVerificationApp.AmqpEventConsumer.StartConsumer(userVerificationApp.Worker)
		if err1 != nil {
			slog.Error("failed to start verification Consumer", err1)
			cancel()
			<-ctx.Done()
		}
	}()

	return userVerificationApp
}
