package v1

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/infra/logger"
)

type UserV1Api struct {
	logger      logger.Logger
	userService inbound.UserUseCase
}

// NewUserApi creates a new UserV1Api structure
func NewUserApi(userService inbound.UserUseCase, log logger.Logger) UserV1Api {
	return UserV1Api{
		logger:      log,
		userService: userService,
	}
}
