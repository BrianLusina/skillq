package users

import "github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"

// UserApi represents a structure for the APO
type UserApi struct {
	// service to handle user management
	userService inbound.UserUseCase
}

// NewApi creates a new User API
func NewApi(userService inbound.UserUseCase) UserApi {
	return UserApi{
		userService: userService,
	}
}
