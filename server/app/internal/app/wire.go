//go:build wireinject
// +build wireinject

package app

import (
	"github.com/BrianLusina/skillq/server/app/di"
	"github.com/BrianLusina/skillq/server/infra/clients/email"
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
	emailConfig email.EmailClientConfig,
) (*App, error) {
	panic(wire.Build(
		New,
		di.LoggerSet,
		di.ProvideUserMongoDbClient,
		di.UserRepositoryAdapterSet,
		di.AmqpClientSet,
		di.AmqpEventPublisherSet,
		di.SendEmailEventPublisherSet,
		di.StoreImageEventPublisherSet,
		di.StorageMinioClientSet,
		di.UserServiceSet,
		di.ProvideUserVerificationMongoDbClient,
		di.UserVerificationRepositoryAdapterSet,
		di.AmqpEventConsumerSet,
		di.UserVerificationServiceSet,
		di.ProvideSendEmailVerificationTaskHandler,
		di.EmailClientSet,
		di.ProvideStoreImageTaskHandler,
	))
}
