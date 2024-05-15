package usersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound/common"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/BrianLusina/skillq/server/utils/security"
	"github.com/BrianLusina/skillq/server/utils/tools"
	"github.com/pkg/errors"
)

// userService is the structure for the business logic handling user management
type userService struct {
	userRepo              repositories.UserRepoPort
	userVerificationRepo  repositories.UserVerificationRepoPort
	eventMessagePublisher amqppublisher.AmqpEventPublisher
	storageClient         storage.StorageClient
}

var _ inbound.UserUseCase = (*userService)(nil)

// New creates a new user service implementation of the user use case
func New(
	userRepo repositories.UserRepoPort,
	userVerificationRepo repositories.UserVerificationRepoPort,
	messagePublisher amqppublisher.AmqpEventPublisher,
	storageClient storage.StorageClient,
) inbound.UserUseCase {
	return &userService{
		userRepo:              userRepo,
		userVerificationRepo:  userVerificationRepo,
		eventMessagePublisher: messagePublisher,
		storageClient:         storageClient,
	}
}

// CreateUser creates a new user in the system
func (svc *userService) CreateUser(ctx context.Context, request inbound.UserRequest) (*inbound.UserResponse, error) {
	// hash password
	hashedPassword, err := security.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := user.New(user.UserParams{
		EntityParams: entity.EntityParams{
			EntityIDParams: entity.EntityIDParams{
				UUID:  id.NewUUID(),
				KeyID: id.NewKeyID(),
				XID:   id.NewXid(),
			},
			EntityTimestampParams: entity.EntityTimestampParams{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Metadata: map[string]any{},
		},
		Name:     request.Name,
		Email:    request.Email,
		Skills:   request.Skills,
		JobTitle: request.JobTitle,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	// create user
	createdUser, err := svc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: move to goroutine
	// create email verification & send it
	if _, err := svc.CreateEmailVerification(ctx, createdUser.UUID()); err != nil {
		return nil, errors.Wrapf(err, "failed to create user email verification")
	}

	// TODO: move to goroutine
	// send image data to persist to blob storage & persist the image URL in the database
	imageUrl, err := svc.UploadUserImage(ctx, createdUser.UUID(), request.Image)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upload user image data")
	}

	// update user image in response
	return mapUserToUserResponse(createdUser.WithImage(imageUrl)), nil
}

// CreateEmailVerification creates a user verification & publishes it to a topic for a listener to send to a user
func (svc *userService) CreateEmailVerification(ctx context.Context, userUUID id.UUID) (user.UserVerification, error) {
	// retrieve the existingUser
	existingUser, err := svc.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to retrieve user %w", err)
	}

	code := security.GenerateCode()

	verificationId := id.NewUUID()
	now := time.Now()

	verification := user.NewVerification(user.UserVerificationParams{
		ID:         verificationId,
		UserId:     existingUser.UUID(),
		Code:       code,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	// create user verification
	if _, err := svc.userVerificationRepo.CreateUserVerification(ctx, verification); err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to create email verification: %w", err)
	}

	event := events.UserEmailVerificationStarted{
		UserUUID: existingUser.UUID(),
		Email:    existingUser.Email(),
		Name:     existingUser.Name(),
		Code:     code,
	}

	eventBytes, err := events.EventToBytes(event)
	if err != nil {
		return user.UserVerification{}, errors.Wrapf(err, "failed to send user email verification")
	}

	// TODO: move to separate goroutine
	if err := svc.eventMessagePublisher.Publish(ctx, eventBytes, "text/plain"); err != nil {
		return user.UserVerification{}, errors.Wrapf(err, "failed to publish user email verified task")
	}

	return verification, nil
}

// GetUserByUUID retrieves a user given their UUID
func (svc *userService) GetUserByUUID(ctx context.Context, userUUID string) (*inbound.UserResponse, error) {
	uuid, err := id.StringToUUID(userUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse user UUID %s", userUUID)
	}

	existingUser, err := svc.userRepo.GetUserByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user %w", err)
	}

	return mapUserToUserResponse(*existingUser), nil
}

func (svc *userService) UploadUserImage(ctx context.Context, userUUID id.UUID, imageData inbound.UserImageRequest) (string, error) {
	url, err := svc.storageClient.Upload(ctx, storage.StorageItem{
		ContentType: imageData.Type,
		Content:     imageData.Content,
		Name:        fmt.Sprintf("%s-image", userUUID),
		Bucket:      fmt.Sprintf("%s-documents", userUUID),
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to store user image")
	}
	return url, nil
}

// GetAllUsers retrieves all users
func (svc *userService) GetAllUsers(ctx context.Context, params common.RequestParams) ([]inbound.UserResponse, error) {
	users, err := svc.userRepo.GetAllUsers(ctx, params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all users")
	}

	return tools.MapWithError(users, func(u user.User, _ int) (inbound.UserResponse, error) {
		return *mapUserToUserResponse(u), nil
	})
}

// GetAllUsersBySkill retrieves all users with a given skill
func (svc *userService) GetAllUsersBySkill(ctx context.Context, skill string, params common.RequestParams) ([]inbound.UserResponse, error) {
	users, err := svc.userRepo.GetAllUsersBySkill(ctx, skill, params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all users with skill %s", skill)
	}

	return tools.MapWithError(users, func(u user.User, _ int) (inbound.UserResponse, error) {
		return *mapUserToUserResponse(u), nil
	})
}

// DeleteUser deletes a given user
func (svc *userService) DeleteUser(ctx context.Context, userId string) error {
	uuid, err := id.StringToUUID(userId)
	if err != nil {
		return errors.Wrapf(err, "failed to parse user ID %s", userId)
	}

	err = svc.userRepo.DeleteUserById(ctx, uuid)
	if err != nil {
		return errors.Wrapf(err, "failed to delete user with ID: %s", userId)
	}

	return nil
}
