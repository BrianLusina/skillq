package taskhandlers

import (
	"context"
	"fmt"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/logger"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/pkg/errors"
)

type storeUserImageTaskHandler struct {
	storageClient                   storage.StorageClient
	userSvc                         inbound.UserService
	emailVerificationEventPublisher amqppublisher.AmqpEventPublisher
	logger                          logger.Logger
}

var _ handlers.EventHandler[tasks.StoreUserImage] = (*storeUserImageTaskHandler)(nil)

func NewStoreImageTasksHandler(
	storageClient storage.StorageClient,
	userSvc inbound.UserService,
	messagePublisher amqppublisher.AmqpEventPublisher,
	logger logger.Logger,
) handlers.EventHandler[tasks.StoreUserImage] {
	return &storeUserImageTaskHandler{
		storageClient:                   storageClient,
		userSvc:                         userSvc,
		emailVerificationEventPublisher: messagePublisher,
		logger:                          logger,
	}
}

func (h *storeUserImageTaskHandler) Handle(ctx context.Context, task tasks.StoreUserImage) error {
	h.logger.Infof("Received tasks store user image, %v", task)

	userUUID, contentType, content, name, bucket := task.UserUUID, task.ContentType, task.Content, task.Name, task.Bucket

	url, err := h.storageClient.Upload(ctx, storage.StorageItem{
		ContentType: contentType,
		Content:     content,
		Name:        name,
		Bucket:      bucket,
	})

	if err != nil {
		return errors.Wrapf(err, "failed to store user image")
	}

	//TODO:  update user data
	user, err := h.userSvc.GetUserByUUID(ctx, userUUID)
	if err != nil {
		msg := fmt.Sprintf("Failed to create email verification for user %s for email %s", contentType.String(), content)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	if err := h.emailVerificationEventPublisher.Publish(ctx, eventBytes, "text/plain"); err != nil {
		msg := fmt.Sprintf("Failed to publish event %v", sendEmailEvent)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, "failed to publish user email sent event")
	}

	return nil
}
