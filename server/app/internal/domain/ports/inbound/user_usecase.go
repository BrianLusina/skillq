package inbound

import (
	"context"
	"time"
)

// UserRequest to create a new user
type UserRequest struct {
	Name     string
	Email    string
	Password string
	Skills   []string
	Image    []byte
	JobTitle string
}

// UserResponse for returning a user
type UserResponse struct {
	UUID      string
	KeyID     string
	XID       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Metadata  map[string]any
	Name      string
	Email     string
	Skills    []string
	ImageUrl  string
	JobTitle  string
}

// UserUserCase contains a method set defining the logic to handle user management in the system
type UserUseCase interface {
	// CreateUser creates a new user in the system
	CreateUser(context.Context, UserRequest) (*UserResponse, error)
}
