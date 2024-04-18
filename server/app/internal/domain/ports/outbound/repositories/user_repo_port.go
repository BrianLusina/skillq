package repositories

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
)

// UserRepoPort handles repository interface
type UserRepoPort interface {
	// CreateUser creates a user in the repository
	CreateUser(context.Context, user.User) (*user.User, error)
}
