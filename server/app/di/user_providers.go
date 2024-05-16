package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/google/wire"
)

var UserServiceSet = wire.NewSet(usersvc.New)

func ProvideUserMongoDbClient(cfg mongodb.MongoDBConfig) mongodb.MongoDBClient[models.UserModel] {
	log := logger.New()
	userMongoDbClient, err := mongodb.New[models.UserModel](cfg, log)
	if err != nil {
		panic(err)
	}
	return userMongoDbClient
}

var UserMongoDbClientSet = wire.NewSet(mongodb.New[models.UserModel])
var UserRepositoryAdapterSet = wire.NewSet(userrepo.New)

var UserVerificationMongoDbClient = wire.NewSet(mongodb.New[models.UserVerificationModel])

func ProvideUserVerificationMongoDbClient(cfg mongodb.MongoDBConfig) mongodb.MongoDBClient[models.UserVerificationModel] {
	log := logger.New()
	userMongoDbClient, err := mongodb.New[models.UserVerificationModel](cfg, log)
	if err != nil {
		panic(err)
	}
	return userMongoDbClient
}

var UserVerificationRepositoryAdapterSet = wire.NewSet(userrepo.NewVerification)