package verificationapp

import (
	"github.com/BrianLusina/skillq/server/app/di"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/google/wire"
)

func InitializeUserVerificationApp(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
) (*UserVerificationApp, func(), error) {
	panic(wire.Build(
		NewUserVerificationApp,
		di.LoggerSet,
		di.UserVerificationRepositoryAdapterSet,
		di.ProvideUserVerificationMongoDbClient,
		di.AmqpClientSet,
		di.AmqpEventPublisherSet,
		di.UserVerificationServiceSet,
		di.UserServiceSet,
	))
}
