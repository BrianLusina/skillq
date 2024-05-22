package inbound

import (
	"context"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// UserRequest to create a new user
type UserRequest struct {
	Name     string
	Email    string
	Password string
	Skills   []string
	Image    UserImageRequest
	JobTitle string
}

// UserImageRequest is the user image data
type UserImageRequest struct {
	// Type is the type of the image
	Type string
	// Content is the content of the image
	Content string
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

// UserService contains a method set defining the logic to handle user management in the system
type UserService interface {
	// CreateUser creates a new user in the system
	CreateUser(context.Context, UserRequest) (*UserResponse, error)

	// GetUserByUUID retrieves a user given their UUID
	GetUserByUUID(context.Context, string) (*UserResponse, error)

	// GetAllUsers retrieves all users
	GetAllUsers(context.Context, common.RequestParams) ([]UserResponse, error)

	// GetAllUsersBySkill retrieves all users with a given skill
	GetAllUsersBySkill(context.Context, string, common.RequestParams) ([]UserResponse, error)

	// UploadUserImage uploads a user image to blob storage & retrieves the image url
	UploadUserImage(context.Context, id.UUID, UserImageRequest) (string, error)

	// UpdateUser updates a user given their ID
	UpdateUser(ctx context.Context, userID string, request UserRequest) (*UserResponse, error)

	// UpdateUserImage updates a user's image
	// UpdateUserImage(ctx context.Context, userID, url string) error

	// DeleteUser deletes a user given their ID
	DeleteUser(context.Context, string) error
}
