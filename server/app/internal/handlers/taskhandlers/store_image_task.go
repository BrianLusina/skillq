package taskhandlers

import (
	"context"
	"fmt"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/logger"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/pkg/errors"
)

type storeUserImageTaskHandler struct {
	storageClient                   storage.StorageClient
	userRepo                        repositories.UserRepoPort
	emailVerificationEventPublisher amqppublisher.AmqpEventPublisher
	logger                          logger.Logger
}

var _ handlers.EventHandler[tasks.StoreUserImage] = (*storeUserImageTaskHandler)(nil)

func NewStoreImageTasksHandler(
	storageClient storage.StorageClient,
	userRepo repositories.UserRepoPort,
	messagePublisher amqppublisher.AmqpEventPublisher,
	logger logger.Logger,
) handlers.EventHandler[tasks.StoreUserImage] {
	return &storeUserImageTaskHandler{
		storageClient:                   storageClient,
		userRepo:                        userRepo,
		emailVerificationEventPublisher: messagePublisher,
		logger:                          logger,
	}
}

func (h *storeUserImageTaskHandler) Handle(ctx context.Context, task tasks.StoreUserImage) error {
	h.logger.Infof("Received tasks store user image, %v", task)

	userID, contentType, content, name, bucket := task.UserUUID, task.ContentType, task.Content, task.Name, task.Bucket
	storageItem := storage.StorageItem{
		ContentType: contentType,
		Content:     content,
		Name:        name,
		Bucket:      bucket,
	}

	url, err := h.storageClient.Upload(ctx, storageItem)
	if err != nil {
		return errors.Wrapf(err, "failed to store user image")
	}
	h.logger.Info("Successfully uploaded user image")

	userUUID, err := id.StringToUUID(userID)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse user ID %s", userID)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	user, err := h.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve user %s", userID)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	updateUser := user.SetImageUrl(url)
	if _, err := h.userRepo.UpdateUser(ctx, *updateUser); err != nil {
		msg := fmt.Sprintf("Failed to update user image %s", userID)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	return nil
}
