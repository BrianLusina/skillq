package v1

import "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"

type UserV1Api struct {
	userService inbound.UserUseCase
}

// NewUserApi creates a new UserV1Api structure
func NewUserApi(userService inbound.UserUseCase) UserV1Api {
	return UserV1Api{
		userService: userService,
	}
}
