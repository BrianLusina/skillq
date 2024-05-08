package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/google/wire"
)

// user provider sets
var UserServiceSet = wire.NewSet(usersvc.New)

func ProvideUserMongoDbClient(cfg mongodb.MongoDBConfig) mongodb.MongoDBClient[models.UserModel] {
	userMongoDbClient, err := mongodb.New[models.UserModel](cfg)
	if err != nil {
		panic(err)
	}
	return userMongoDbClient
}

var UserMongoDbClientSet = wire.NewSet(mongodb.New[models.UserModel])
var UserRepositoryAdapterSet = wire.NewSet(userrepo.New)

var AmqpClientSet = wire.NewSet(amqp.NewAmqpClient)
var EventPublisherSet = wire.NewSet(publisher.NewPublisher)

var UserVerificationMongoDbClient = wire.NewSet(mongodb.New[models.UserVerificationModel])

func ProvideUserVerificationMongoDbClient(cfg mongodb.MongoDBConfig) mongodb.MongoDBClient[models.UserVerificationModel] {
	userMongoDbClient, err := mongodb.New[models.UserVerificationModel](cfg)
	if err != nil {
		panic(err)
	}
	return userMongoDbClient
}

var UserVerificationRepositoryAdapterSet = wire.NewSet(userrepo.NewVerification)

var UserAppLoggerSet = wire.NewSet(logger.New)

var UserStorageClientSet = wire.NewSet(minio.NewClient)
