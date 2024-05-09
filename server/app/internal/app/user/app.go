package userapp

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
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

	UserVerificationRepo repositories.UserVerificationRepoPort
	UserRepo             repositories.UserRepoPort

	AmqpClient         *amqp.AmqpClient
	AmqpEventPublisher amqppublisher.AmqpEventPublisher
	StorageClient      storage.StorageClient

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
	amqpClient *amqp.AmqpClient,
	amqpEventPublisher amqppublisher.AmqpEventPublisher,
	storageClient storage.StorageClient,
	userRepo repositories.UserRepoPort,
	userVerificationRepo repositories.UserVerificationRepoPort,
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
		AmqpEventPublisher:            amqpEventPublisher,
		StorageClient:                 storageClient,
		UserVerificationRepo:          userVerificationRepo,
		UserRepo:                      userRepo,
		UserSvc:                       userSvc,
	}
}
