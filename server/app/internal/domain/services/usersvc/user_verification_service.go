package usersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/utils/security"
	"github.com/pkg/errors"
)

// userVerificationService is the structure for the business logic handling user verification
type userVerificationService struct {
	userSvc              inbound.UserService
	userVerificationRepo repositories.UserVerificationRepoPort
}

var _ inbound.UserVerificationService = (*userVerificationService)(nil)

// NewVerification creates a new user service implementation of the user use case
func NewVerification(
	userSvc inbound.UserService,
	userVerificationRepo repositories.UserVerificationRepoPort,
) inbound.UserVerificationService {
	return &userVerificationService{
		userSvc:              userSvc,
		userVerificationRepo: userVerificationRepo,
	}
}

// CreateEmailVerification creates a user verification & publishes it to a topic for a listener to send to a user
func (svc *userVerificationService) CreateEmailVerification(ctx context.Context, userUUID string, email string) (user.UserVerification, error) {
	// retrieve the existingUser
	if _, err := svc.userSvc.GetUserByUUID(ctx, userUUID); err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to retrieve user %w", err)
	}

	code := security.GenerateCode()

	uuid, err := id.StringToUUID(userUUID)
	if err != nil {
		return user.UserVerification{}, errors.Wrapf(err, "failed to parse user UUID %s", userUUID)
	}

	verificationId := id.NewUUID()
	now := time.Now()

	verification := user.NewVerification(user.UserVerificationParams{
		ID:         verificationId,
		UserId:     uuid,
		Code:       code,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	// create user verification
	if _, err := svc.userVerificationRepo.CreateUserVerification(ctx, verification); err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to create user verification: %w", err)
	}

	return verification, nil
}

// VerifyEmail verifies a user email
func (svc *userVerificationService) VerifyEmail(ctx context.Context, request inbound.VerifyEmailRequest) error {
	userId, code := request.UserID, request.Code
	if _, err := svc.userSvc.GetUserByUUID(ctx, userId); err != nil {
		return fmt.Errorf("failed to retrieve user %w", err)
	}
	userUUID, err := id.StringToUUID(userId)
	if err != nil {
		return errors.Wrapf(err, "failed to parse user UUID %s", userUUID)
	}

	verification, err := svc.userVerificationRepo.GetUserVerificationByUUID(ctx, userUUID)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve user verification for code %s", code)
	}

	if verification.Code() != code {
		return fmt.Errorf("invalid code provided %s", code)
	}

	err = svc.userVerificationRepo.UpdateUserVerification(ctx, repositories.UpdateUserVerificationRequest{
		UserID:     userUUID,
		IsVerified: true,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update user's verification status")
	}

	return nil
}
