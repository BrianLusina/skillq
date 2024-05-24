package userapp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqpconsumer "github.com/BrianLusina/skillq/server/infra/messaging/amqp/consumer"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type (
	// UserApp is a structure for the user application
	UserApp struct {
		MongoDbConfig mongodb.MongoDBConfig
		AmqpConfig    amqp.Config
		MinioConfig   minio.Config

		Logger logger.Logger

		UsersMongoDbClient mongodb.MongoDBClient[models.UserModel]

		UserRepo repositories.UserRepoPort

		AmqpClient         *amqp.AmqpClient
		AmqpEventPublisher amqppublisher.AmqpEventPublisher
		AmqpEventConsumer  amqpconsumer.AmqpEventConsumer
		StorageClient      storage.StorageClient

		UserSvc inbound.UserService

		EmailVerificationSentHandler    handlers.EventHandler[events.EmailVerificationSent]
		EmailVerificationStartedHandler handlers.EventHandler[events.EmailVerificationStarted]
		StoreImageTaskHandler           handlers.EventHandler[tasks.StoreUserImage]
	}
)

// New creates a new UserApp
func New(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
	logger logger.Logger,
	usersMongoDbClient mongodb.MongoDBClient[models.UserModel],
	amqpClient *amqp.AmqpClient,
	amqpEventPublisher amqppublisher.AmqpEventPublisher,
	amqpEventConsumer amqpconsumer.AmqpEventConsumer,
	storageClient storage.StorageClient,
	userRepo repositories.UserRepoPort,
	userSvc inbound.UserService,

	emailVerificationSentHandler handlers.EventHandler[events.EmailVerificationSent],
	emailVerificationStartedHandler handlers.EventHandler[events.EmailVerificationStarted],

	storeImageTaskHandler handlers.EventHandler[tasks.StoreUserImage],
) *UserApp {
	return &UserApp{
		MongoDbConfig:      mongodbConfig,
		AmqpConfig:         amqpConfig,
		MinioConfig:        minioConfig,
		Logger:             logger,
		UsersMongoDbClient: usersMongoDbClient,
		AmqpClient:         amqpClient,
		AmqpEventPublisher: amqpEventPublisher,
		AmqpEventConsumer:  amqpEventConsumer,
		StorageClient:      storageClient,
		UserRepo:           userRepo,
		UserSvc:            userSvc,

		EmailVerificationSentHandler:    emailVerificationSentHandler,
		EmailVerificationStartedHandler: emailVerificationStartedHandler,
		StoreImageTaskHandler:           storeImageTaskHandler,
	}
}

func (app *UserApp) Worker(ctx context.Context, messages <-chan any) {
	for message := range messages {
		delivery := message.(rabbitmq.Delivery)
		app.Logger.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		app.Logger.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "email-verification-started":
			var payload events.EmailVerificationStarted

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				app.Logger.Error("failed to Unmarshal message", err)
			}

			err = app.EmailVerificationStartedHandler.Handle(ctx, &payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					app.Logger.Error("failed to delivery.Reject", err)
				}

				app.Logger.Error("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					app.Logger.Error("failed to acknowledge delivery", err)
				}
			}

		case "email-verification-sent":
			var payload events.EmailVerificationSent

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = app.EmailVerificationSentHandler.Handle(ctx, &payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					slog.Error("failed to delivery.Reject", err)
				}

				slog.Error("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					slog.Error("failed to acknowledge delivery", err)
				}
			}

		default:
			app.Logger.Info("default")
		}
	}
}
