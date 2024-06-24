package userv1

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/infra/logger"
)

type UserV1Api struct {
	logger                  logger.Logger
	userService             inbound.UserService
	userVerificationService inbound.UserVerificationService
}

// NewUserApi creates a new UserV1Api structure
func NewUserApi(userService inbound.UserService, userVerificationService inbound.UserVerificationService, log logger.Logger) UserV1Api {
	return UserV1Api{
		logger:                  log,
		userService:             userService,
		userVerificationService: userVerificationService,
	}
}
