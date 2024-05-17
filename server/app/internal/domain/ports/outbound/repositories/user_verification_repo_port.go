package repositories

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// UpdateUserVerificationRequest represents the required fields to update a user verification status
type UpdateUserVerificationRequest struct {
	UserID     id.UUID
	IsVerified bool
}

// UserVerificationRepoPort handles user verification repository interface
type UserVerificationRepoPort interface {
	// CreateUserVerification creates a user verification in the repository
	CreateUserVerification(context.Context, user.UserVerification) (*user.UserVerification, error)

	// GetUserVerificationByUUID retrieves a user verification given the UUID
	GetUserVerificationByUUID(context.Context, id.UUID) (*user.UserVerification, error)

	// GetUserVerificationByCode retrieves a user verification given the code
	GetUserVerificationByCode(context.Context, string) (*user.UserVerification, error)

	// UpdateUserVerification updates the user verification
	UpdateUserVerification(context.Context, UpdateUserVerificationRequest) error
}
