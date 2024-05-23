package userapp

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
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
		StorageClient:      storageClient,
		UserRepo:           userRepo,
		UserSvc:            userSvc,

		EmailVerificationSentHandler:    emailVerificationSentHandler,
		EmailVerificationStartedHandler: emailVerificationStartedHandler,
		StoreImageTaskHandler:           storeImageTaskHandler,
	}
}
