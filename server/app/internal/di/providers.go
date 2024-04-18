package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/services/usersvc"
	"github.com/google/wire"
)

var UserServiceSet = wire.NewSet(usersvc.New)
