//go:build wireinject
// +build wireinject

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
) (*UserVerificationApp, error) {
	panic(wire.Build(
		NewUserVerificationApp,
		di.LoggerSet,
		di.UserVerificationRepositoryAdapterSet,
		di.ProvideUserVerificationMongoDbClient,
		di.AmqpClientSet,
		di.AmqpEventPublisherSet,
		di.AmqpEventConsumerSet,
		di.UserVerificationServiceSet,
		di.ProvideUserMongoDbClient,
		di.UserRepositoryAdapterSet,
		di.StorageMinioClientSet,
		di.ProvideStoreImageTaskHandler,
		di.UserServiceSet,
	))
}
