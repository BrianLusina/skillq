package verificationapp

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

type (
	UserVerificationApp struct {
		MongoDbConfig mongodb.MongoDBConfig
		AmqpConfig    amqp.Config
		MinioConfig   minio.Config

		Logger logger.Logger

		AmqpClient         *amqp.AmqpClient
		AmqpEventPublisher amqppublisher.AmqpEventPublisher
		StorageClient      storage.StorageClient

		UserVerificationMongoDbClient mongodb.MongoDBClient[models.UserVerificationModel]
		UserVerificationRepo          repositories.UserVerificationRepoPort

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
		UserVerificationRepo:          userVerificationRepo,
		UserVerificationMongoDbClient: userVerificationMongoDbClient,
		UserVerificationSvc:           userVerificationService,
	}
}
