// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package verificationapp

import (
	"github.com/BrianLusina/skillq/server/app/di"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userverification"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
)

// Injectors from wire.go:

func InitializeUserVerificationApp(mongodbConfig mongodb.MongoDBConfig, amqpConfig amqp.Config, minioConfig minio.Config) (*UserVerificationApp, error) {
	loggerLogger := logger.New()
	mongoDBClient := di.ProvideUserMongoDbClient(mongodbConfig)
	amqpClient, err := amqp.NewAmqpClient(amqpConfig, loggerLogger)
	if err != nil {
		return nil, err
	}
	amqpEventPublisher, err := amqppublisher.NewPublisher(amqpClient, loggerLogger)
	if err != nil {
		return nil, err
	}
	storageClient, err := minio.NewClient(minioConfig, loggerLogger)
	if err != nil {
		return nil, err
	}
	userRepoPort := userrepo.New(mongoDBClient)
	eventHandler := di.ProvideStoreImageTaskHandler(storageClient, userRepoPort, amqpEventPublisher)
	mongodbMongoDBClient := di.ProvideUserVerificationMongoDbClient(mongodbConfig)
	userVerificationRepoPort := userverificationrepo.New(mongodbMongoDBClient)
	userService := usersvc.New(userRepoPort, amqpEventPublisher, storageClient)
	userVerificationService := usersvc.NewVerification(userService, userVerificationRepoPort, amqpEventPublisher)
	userVerificationApp := NewUserVerificationApp(mongodbConfig, amqpConfig, minioConfig, loggerLogger, mongoDBClient, amqpClient, amqpEventPublisher, eventHandler, mongodbMongoDBClient, userVerificationRepoPort, userVerificationService)
	return userVerificationApp, nil
}
