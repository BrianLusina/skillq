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
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
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

	cleanup := prepareApp(ctx, cancel, app, cfg)

	// Start the server
	err = app.Listen(fmt.Sprintf(":%d", cfg.HTTP.Port))
	appLogger.Fatalf("Failed to start application: %v", err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		cleanup()
		appLogger.Info("signal.Notify", v)
	case done := <-ctx.Done():
		cleanup()
		appLogger.Info("ctx.Done", "app done", done)
	}
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, app *fiber.App, cfg *config.Config) func() {
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
	// 		CollectionName: cfg.MongoDB.Collections["userVerification"].Name,
	// 	},
	// }

	// collections := cfg.MongoDB.Collections

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

	userApp, userAppCleanup := prepareUserApp(ctx, cancel, userMongodbConfig, amqpConfig, minioConfig)
	appLogger := logger.New()

	// routing
	userApi := userv1.NewUserApi(userApp.UserSvc, appLogger)

	userApi.RegisterHandlers(app)

	return userAppCleanup
}

func prepareUserApp(ctx context.Context, cancel context.CancelFunc, mongoDbConfig mongodb.MongoDBConfig, amqpConfig amqp.Config, minioConfig minio.Config) (*userapp.UserApp, func()) {
	userApp, cleanup, err := userapp.InitializeUserApp(mongoDbConfig, amqpConfig, minioConfig)
	if err != nil {
		slog.Error("failed init user app", err)
		cancel()
		<-ctx.Done()
	}

	// Configure publisher and start workers
	// userApp.AmqpEventPublisher.Configure()

	// a.BaristaOrderPub.Configure(
	// 	pkgPublisher.ExchangeName("barista-order-exchange"),
	// 	pkgPublisher.BindingKey("barista-order-routing-key"),
	// 	pkgPublisher.MessageTypeName("barista-order-created"),
	// )

	// a.KitchenOrderPub.Configure(
	// 	pkgPublisher.ExchangeName("kitchen-order-exchange"),
	// 	pkgPublisher.BindingKey("kitchen-order-routing-key"),
	// 	pkgPublisher.MessageTypeName("kitchen-order-created"),
	// )

	// a.Consumer.Configure(
	// 	pkgConsumer.ExchangeName("counter-order-exchange"),
	// 	pkgConsumer.QueueName("counter-order-queue"),
	// 	pkgConsumer.BindingKey("counter-order-routing-key"),
	// 	pkgConsumer.ConsumerTag("counter-order-consumer"),
	// )

	// go func() {
	// 	err1 := a.Consumer.StartConsumer(a.Worker)
	// 	if err1 != nil {
	// 		slog.Error("failed to start Consumer", err1)
	// 		cancel()
	// 		<-ctx.Done()
	// 	}
	// }()

	return userApp, cleanup
}
