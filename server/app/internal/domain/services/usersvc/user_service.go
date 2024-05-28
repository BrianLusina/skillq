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
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/domain/entity"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/BrianLusina/skillq/server/utils/security"
	"github.com/BrianLusina/skillq/server/utils/tools"
	"github.com/pkg/errors"
)

// userService is the structure for the business logic handling user management
type userService struct {
	userRepo         repositories.UserRepoPort
	messagePublisher amqppublisher.AmqpEventPublisher
	storageClient    storage.StorageClient
}

var _ inbound.UserService = (*userService)(nil)

// New creates a new user service implementation of the user use case
func New(
	userRepo repositories.UserRepoPort,
	messagePublisher amqppublisher.AmqpEventPublisher,
	storageClient storage.StorageClient,
) inbound.UserService {
	return &userService{
		userRepo:         userRepo,
		messagePublisher: messagePublisher,
		storageClient:    storageClient,
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

	event := events.EmailVerificationStarted{
		UserUUID: createdUser.UUID(),
		Email:    createdUser.Email(),
		Name:     createdUser.Name(),
	}

	message := messaging.Message{
		Topic:       event.Identity(),
		Payload:     event,
		ContentType: "text/plain",
	}

	// publish event
	if err := svc.messagePublisher.Publish(ctx, message); err != nil {
		return nil, errors.Wrapf(err, "failed to publish user email verification event: %v", message)
	}

	storeUserImageTask := tasks.StoreUserImageTask{
		ContentType: request.Image.Type,
		Content:     request.Image.Content,
		Name:        fmt.Sprintf("%s-image", createdUser.UUID()),
		Bucket:      fmt.Sprintf("%s-documents", createdUser.UUID()),
	}

	storeImageMessage := messaging.Message{
		Topic:       storeUserImageTask.Identity(),
		Payload:     storeUserImageTask,
		ContentType: "text/plain",
	}

	// publish store image task
	if err := svc.messagePublisher.Publish(ctx, storeImageMessage); err != nil {
		return nil, errors.Wrapf(err, "failed to publish store user image task: %v", storeImageMessage)
	}

	// update user image in response
	return mapUserToUserResponse(*createdUser), nil
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

// UpdateUser updates a user given their ID
func (svc *userService) UpdateUser(ctx context.Context, userID string, request inbound.UserRequest) (*inbound.UserResponse, error) {
	userUUID, err := id.StringToUUID(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse user ID %s", userUUID)
	}

	// get user
	existingUser, err := svc.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user %w", err)
	}

	// update fields
	url, err := svc.UploadUserImage(ctx, userUUID, request.Image)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upload user image %s", request.Image)
	}

	existingUpdatedUser, err := existingUser.SetName(request.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update user name %s", request.Name)
	}

	existingUpdatedUser, err = existingUpdatedUser.SetEmail(request.Email)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update user email %s", request.Email)
	}

	existingUpdatedUser, err = existingUpdatedUser.SetJobTitle(request.JobTitle)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update user job title %s", request.JobTitle)
	}

	existingUpdatedUser = existingUpdatedUser.SetSkills(request.Skills).SetImageUrl(url)

	updatedUser, err := svc.userRepo.UpdateUser(ctx, *existingUpdatedUser)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update user %s", userID)
	}

	return mapUserToUserResponse(*updatedUser), nil
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
