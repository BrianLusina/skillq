package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/models"
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/BrianLusina/skillq/server/infra/mongodb"
	"github.com/google/wire"
)

// user provider sets
var UserServiceSet = wire.NewSet(usersvc.New)
var UserMongoDbClient = wire.NewSet(mongodb.New[models.UserModel])
var UserRepositoryAdapterSet = wire.NewSet(userrepo.New)

var UserVerificationMongoDbClient = wire.NewSet(mongodb.New[models.UserVerificationModel])
var UserVerificationRepositoryAdapterSet = wire.NewSet(userrepo.NewVerification)
