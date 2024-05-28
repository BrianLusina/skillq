package verificationapp

import (
	"context"
	"encoding/json"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
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
	UserVerificationApp struct {
		MongoDbConfig mongodb.MongoDBConfig
		AmqpConfig    amqp.Config
		MinioConfig   minio.Config

		Logger logger.Logger

		AmqpClient         *amqp.AmqpClient
		AmqpEventPublisher amqppublisher.AmqpEventPublisher
		AmqpEventConsumer  amqpconsumer.AmqpEventConsumer
		StorageClient      storage.StorageClient

		UserVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel]
		UserVerificationRepo          repositories.UserVerificationRepoPort

		StoreImageTaskHandler handlers.EventHandler[tasks.StoreUserImageTask]

		UserVerificationSvc inbound.UserVerificationService
	}
)

// NewUserVerificationApp creates a new user verification app
func NewUserVerificationApp(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
	logger logger.Logger,
	usersMongoDbClient mongodb.MongoDBClient[models.UserModel],
	amqpClient *amqp.AmqpClient,
	amqpEventPublisher amqppublisher.AmqpEventPublisher,
	amqpEventConsumer amqpconsumer.AmqpEventConsumer,
	storeImageTaskHandler handlers.EventHandler[tasks.StoreUserImageTask],

	userVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel],
	userVerificationRepo repositories.UserVerificationRepoPort,
	userVerificationService inbound.UserVerificationService,
) *UserVerificationApp {
	return &UserVerificationApp{
		MongoDbConfig:                 mongodbConfig,
		AmqpConfig:                    amqpConfig,
		MinioConfig:                   minioConfig,
		Logger:                        logger,
		AmqpClient:                    amqpClient,
		AmqpEventPublisher:            amqpEventPublisher,
		AmqpEventConsumer:             amqpEventConsumer,
		UserVerificationRepo:          userVerificationRepo,
		UserVerificationMongoDbClient: userVerificationMongoDbClient,
		UserVerificationSvc:           userVerificationService,
		StoreImageTaskHandler:         storeImageTaskHandler,
	}
}

func (app *UserVerificationApp) Worker(ctx context.Context, messages <-chan rabbitmq.Delivery) {
	for message := range messages {
		app.Logger.Info("processDeliveries", "delivery_tag", message.DeliveryTag)
		app.Logger.Info("received", "delivery_type", message.Type)

		switch message.Type {
		case string(tasks.StoreUserImageTaskName):
			var payload tasks.StoreUserImageTask

			err := json.Unmarshal(message.Body, &payload)
			if err != nil {
				app.Logger.Error("failed to Unmarshal message", err)
			}

			err = app.StoreImageTaskHandler.Handle(ctx, &payload)

			if err != nil {
				if err = message.Reject(false); err != nil {
					app.Logger.Error("failed to delivery.Reject", err)
				}

				app.Logger.Error("failed to process delivery", err)
			} else {
				err = message.Ack(false)
				if err != nil {
					app.Logger.Error("failed to acknowledge delivery", err)
				}
			}

		default:
			app.Logger.Info("default")
		}
	}

}
