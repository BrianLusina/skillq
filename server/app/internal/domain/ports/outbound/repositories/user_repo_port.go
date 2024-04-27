package repositories

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// UserRepoPort handles repository interface
type UserRepoPort interface {
	// CreateUser creates a user in the repository
	CreateUser(context.Context, user.User) (*user.User, error)

	// GetUserByUUID retrieves a user given their UUID
	GetUserByUUID(context.Context, id.UUID) (*user.User, error)

	// GetAllUsers retrieves all users
	GetAllUsers(context.Context) ([]user.User, error)
}
