package main

import (
	"context"
	"log"
	"log/slog"

	userv1 "github.com/BrianLusina/skillq/server/app/api/rest/routes/users/v1"
	"github.com/BrianLusina/skillq/server/app/cmd/config"
	userapp "github.com/BrianLusina/skillq/server/app/internal/app/user"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed set max procs", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", err)
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// TODO: setup config
	app := fiber.New()

	go func() {
		defer app.Server().Shutdown()
		<-ctx.Done()
	}()

	// middleware
	app.Use(cors.New())

	cleanup := prepareApp(ctx, cancel, cfg)

	appLogger := logger.New()

	//configuration

	var MONGO_URL = "<your_connection_string>"

	// routing

	usersMongoDbClient, err := mongodb.New[any](mongodb.MongoDBConfig{})
	if err != nil {
		log.Fatalf("Failed to startup application: Err: %v", err)
	}

	userRepo := userrepo.New(usersMongoDbClient)

	userVerificationMongoDbClient, err := mongodb.New[any](mongodb.MongoDBConfig{})
	userVerificationRepo := userrepo.NewVerification(userVerificationMongoDbClient)

	userService := usersvc.New(userRepo, userVerificationRepo)

	userApi := userv1.NewUserApi(userService, appLogger)

	userApi.RegisterHandlers(app)

	// Start the server
	err = app.Listen(":3000")
	log.Fatalf("Failed to start application: %v", err)
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) func() {
	userAppCleanup := prepareUserApp(ctx, cancel, cfg)

	return userAppCleanup
}

func prepareUserApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) func() {
	mongoDbConfig := mongodb.MongoDBConfig{
		Client: mongodb.ClientOptions{
			Host:        cfg.MongoDB.Host,
			Port:        cfg.MongoDB.Port,
			User:        cfg.MongoDB.User,
			Password:    cfg.MongoDB.Password,
			RetryWrites: cfg.MongoDB.RetryWrites,
		},
		DBConfig: mongodb.DatabaseConfig{
			DatabaseName:   cfg.MongoDB.DatabaseName,
			CollectionName: cfg.MongoDB.CollectionName,
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

	a, cleanup, err := userapp.InitializeUserApp(mongoDbConfig, amqpConfig, minioConfig)
	if err != nil {
		slog.Error("failed init user app", err)
		cancel()
		<-ctx.Done()
	}

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

	return cleanup
}
