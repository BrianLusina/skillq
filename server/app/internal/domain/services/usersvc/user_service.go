package usersvc

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/BrianLusina/skillq/server/utils/security"
)

// userService is the structure for the business logic handling user management
type userService struct {
	userRepo         repositories.UserRepoPort
	messagePublisher messaging.Publisher
}

var _ inbound.UserUseCase = (*userService)(nil)

// New creates a new user service implementation of the user use case
func New(userRepo repositories.UserRepoPort, messagePublisher messaging.Publisher) inbound.UserUseCase {
	return &userService{
		userRepo:         userRepo,
		messagePublisher: messagePublisher,
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
		Name:      request.Name,
		Email:     request.Email,
		ImageData: request.Image,
		Skills:    request.Skills,
		JobTitle:  request.JobTitle,
		Password:  hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	// create user
	u, err := svc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// send verification email task

	// send image task to persist to blob storage

	return mapUserToUserResponse(*u), nil
}

// CreateEmailVerification creates a user verification & publishes it to a topic for a listener to send to a user
func (svc *userService) CreateEmailVerification(ctx context.Context, userUUID id.UUID) (user.UserVerification, error) {
	// retrieve the existingUser
	existingUser, err := svc.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to retrieve user %w", err)
	}

	// TODO: move this to utility that handles code generation
	// generate 4 digit code
	randNum := rand.Intn(10000)
	code := fmt.Sprintf("%04d", randNum)

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
	if _, err := svc.userRepo.CreateUserVerification(ctx, verification); err != nil {
		return user.UserVerification{}, fmt.Errorf("failed to create email verification: %w", err)
	}

	// TODO: create event types that will be used to publish events
	// publish to queue/topic
	// run in separate goroutine to avoid blocking calls on this execution path
	go svc.messagePublisher.Publish(ctx, "", "", []byte{})

	return verification, nil
}

// GetUserByUUID retrieves a user given their UUID
func (svc *userService) GetUserByUUID(ctx context.Context, userUUID id.UUID) (*inbound.UserResponse, error) {
	panic("not yet implemented")
}
