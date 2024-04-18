package usersvc

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/app/internal/domain/entities/user"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/utils/security"
)

// userService is the structure for the business logic handling user management
type userService struct {
	userRepo repositories.UserRepoPort
}

var _ inbound.UserUseCase = (*userService)(nil)

// New creates a new user service implementation of the user use case
func New(userRepo repositories.UserRepoPort) inbound.UserUseCase {
	return &userService{
		userRepo: userRepo,
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
