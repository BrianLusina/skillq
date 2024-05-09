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

// InitializeUserApp initializes the user application
func InitializeUserApp(
	mongodbConfig mongodb.MongoDBConfig,
	amqpConfig amqp.Config,
	minioConfig minio.Config,
) (*UserApp, func(), error) {
	panic(wire.Build(
		New,
		di.LoggerSet,
		di.ProvideUserMongoDbClient,
		di.ProvideUserVerificationMongoDbClient,
		di.UserVerificationRepositoryAdapterSet,
		di.UserRepositoryAdapterSet,
		di.AmqpClientSet,
		di.AmqpEventPublisherSet,
		di.StorageMinioClientSet,
		di.UserServiceSet,
	))
}
