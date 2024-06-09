package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/database/repositories/userrepo"
	userverificationrepo "github.com/BrianLusina/skillq/server/app/internal/database/repositories/userverification"
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/google/wire"
)

var UserServiceSet = wire.NewSet(usersvc.New)
var UserRepositoryAdapterSet = wire.NewSet(userrepo.New)

var UserVerificationServiceSet = wire.NewSet(usersvc.NewVerification)
var UserVerificationRepositoryAdapterSet = wire.NewSet(userverificationrepo.New)
