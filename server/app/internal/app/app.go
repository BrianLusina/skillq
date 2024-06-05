package app

import (
	"context"
	"encoding/json"

	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/clients/email"
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
	// App is a structure for the user application
	App struct {
		MongoDbConfig mongodb.MongoDBConfig
		AmqpConfig    amqp.Config
		MinioConfig   minio.Config

		Logger logger.Logger

		UsersMongoDbClient mongodb.MongoDBClient[models.UserModel]

		UserRepo repositories.UserRepoPort

		AmqpClient         *amqp.AmqpClient
		AmqpEventPublisher amqppublisher.AmqpEventPublisher
		AmqpEventConsumer  amqpconsumer.AmqpEventConsumer

		SendEmailEventPublisher  publishers.SendEmailEventPublisherPort
		StoreImageEventPublisher publishers.StoreImageEventPublisherPort

		StorageClient storage.StorageClient

		UserSvc inbound.UserService

		UserVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel]
		UserVerificationRepo          repositories.UserVerificationRepoPort
		UserVerificationSvc           inbound.UserVerificationService

		SendEmailVerificationTaskHandler handlers.EventHandler[tasks.SendEmailVerification]
		StoreImageTaskHandler            handlers.EventHandler[tasks.StoreUserImageTask]

		EmailClient email.EmailClient
	}
)

// New creates a new UserApp
func New(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
	emailConfig email.EmailClientConfig,

	logger logger.Logger,

	amqpClient *amqp.AmqpClient,
	amqpEventPublisher amqppublisher.AmqpEventPublisher,
	amqpEventConsumer amqpconsumer.AmqpEventConsumer,

	sendEmailEventPublisher publishers.SendEmailEventPublisherPort,
	storeImageEventPublisher publishers.StoreImageEventPublisherPort,

	storageClient storage.StorageClient,

	userRepo repositories.UserRepoPort,
	usersMongoDbClient mongodb.MongoDBClient[models.UserModel],
	userSvc inbound.UserService,

	userVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel],
	userVerificationRepo repositories.UserVerificationRepoPort,
	userVerificationService inbound.UserVerificationService,

	sendEmailVerificationHandler handlers.EventHandler[tasks.SendEmailVerification],

	storeImageTaskHandler handlers.EventHandler[tasks.StoreUserImageTask],

	emailClient email.EmailClient,
) *App {
	return &App{
		MongoDbConfig:      mongodbConfig,
		AmqpConfig:         amqpConfig,
		MinioConfig:        minioConfig,
		Logger:             logger,
		UsersMongoDbClient: usersMongoDbClient,

		AmqpClient:         amqpClient,
		AmqpEventPublisher: amqpEventPublisher,
		AmqpEventConsumer:  amqpEventConsumer,

		SendEmailEventPublisher:  sendEmailEventPublisher,
		StoreImageEventPublisher: storeImageEventPublisher,

		StorageClient: storageClient,

		UserRepo: userRepo,
		UserSvc:  userSvc,

		UserVerificationRepo:          userVerificationRepo,
		UserVerificationMongoDbClient: userVerificationMongoDbClient,
		UserVerificationSvc:           userVerificationService,

		SendEmailVerificationTaskHandler: sendEmailVerificationHandler,
		StoreImageTaskHandler:            storeImageTaskHandler,

		EmailClient: emailClient,
	}
}

func (app *App) Worker(ctx context.Context, messages <-chan rabbitmq.Delivery) {
	for message := range messages {
		app.Logger.Infof("Processing message with Tag: %s & Type: %s", message.DeliveryTag, message.Type)

		switch message.Type {
		case string(events.EmailVerificationStartedName):
			var payload tasks.SendEmailVerification

			err := json.Unmarshal(message.Body, &payload)
			if err != nil {
				app.Logger.Error("failed to Unmarshal message", err)
			}

			err = app.SendEmailVerificationTaskHandler.Handle(ctx, &payload)

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
