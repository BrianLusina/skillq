package inbound

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
)

// VerifyEmailRequest defines the structure for verifying a user's email.
type VerifyEmailRequest struct {
	Code   string
	UserID string
}

// UserVerificationService contains a method set defining the logic to handle user verification in the system
type UserVerificationService interface {

	// CreateEmailVerification creates user verification structure that is used to verify a user's email address
	CreateEmailVerification(ctx context.Context, userUUID string, email string) (user.UserVerification, error)

	// VerifyEmail verifies a user email
	VerifyEmail(context.Context, VerifyEmailRequest) error
}
