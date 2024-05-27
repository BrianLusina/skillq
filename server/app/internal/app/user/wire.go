//go:build wireinject
// +build wireinject

package userapp

import (
	"github.com/BrianLusina/skillq/server/app/di"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/google/wire"
)

// InitApp initializes the user application
func InitApp(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
) (*UserApp, error) {
	panic(wire.Build(
		New,
		di.LoggerSet,
		di.ProvideUserMongoDbClient,
		di.UserRepositoryAdapterSet,
		di.AmqpClientSet,
		di.AmqpEventPublisherSet,
		di.StorageMinioClientSet,
		di.UserServiceSet,
		di.ProvideUserVerificationMongoDbClient,
		di.UserVerificationRepositoryAdapterSet,
		di.AmqpEventConsumerSet,
		di.UserVerificationServiceSet,
		di.ProvideEmailVerificationSentEventHandler,
		di.ProvideEmailVerificationStartedEventHandler,
		di.ProvideStoreImageTaskHandler,
	))
}
