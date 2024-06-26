package repositories

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// UserRepoPort handles repository interface
type UserRepoPort interface {
	// CreateUser creates a user in the repository
	CreateUser(context.Context, user.User) (*user.User, error)

	// GetUserByUUID retrieves a user given their UUID
	GetUserByUUID(context.Context, id.UUID) (*user.User, error)

	// GetAllUsers retrieves all users
	GetAllUsers(context.Context, common.RequestParams) ([]user.User, error)

	// GetAllUsersBySkill retrieves all the users of a given skill
	GetAllUsersBySkill(ctx context.Context, skill string, params common.RequestParams) ([]user.User, error)

	// UpdateUser updates a user given their ID
	UpdateUser(context.Context, user.User) (*user.User, error)

	// DeleteUserById deletes a given user by the ID
	DeleteUserById(ctx context.Context, userID id.UUID) error
}
