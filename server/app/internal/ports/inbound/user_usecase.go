package inbound

import (
	"context"
)

// UserRequest to create a new user
type UserRequest struct {
}

// UserResponse for returning a user
type UserResponse struct {
}

// UserUserCase contains a method set defining the logic to handle user management in the system
type UserUseCase interface {
	// CreateUser creates a new user in the system
	CreateUser(context.Context, UserRequest) (UserResponse, error)
}
