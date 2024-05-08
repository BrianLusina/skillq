package userapp

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
)

// UserApp is a structure for the user application
type UserApp struct {
	MongoDbConfig mongodb.MongoDBConfig
	AmqpConfig    amqp.Config
	MinioConfig   minio.Config

	Logger logger.Logger

	UsersMongoDbClient            mongodb.MongoDBClient[models.UserModel]
	UserVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel]

	AmqpClient     amqp.AmqpClient
	EventPublisher messaging.Publisher
	StorageClient  storage.StorageClient

	UserSvc inbound.UserUseCase
}

// New creates a new UserApp
func New(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
	logger logger.Logger,
	usersMongoDbClient mongodb.MongoDBClient[models.UserModel],
	userVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel],
	amqpClient amqp.AmqpClient,
	eventPublisher messaging.Publisher,
	storageClient storage.StorageClient,
	userSvc inbound.UserUseCase,
) *UserApp {
	return &UserApp{
		MongoDbConfig:                 mongodbConfig,
		AmqpConfig:                    amqpConfig,
		MinioConfig:                   minioConfig,
		Logger:                        logger,
		UsersMongoDbClient:            usersMongoDbClient,
		UserVerificationMongoDbClient: userVerificationMongoDbClient,
		AmqpClient:                    amqpClient,
		EventPublisher:                eventPublisher,
		StorageClient:                 storageClient,
		UserSvc:                       userSvc,
	}
}
