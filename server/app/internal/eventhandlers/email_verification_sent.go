package eventhandlers

import (
	"context"
	"fmt"

	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/pkg/errors"
)

type emailVerificationSentEventHandler struct {
	userVerificationSvc             inbound.UserVerificationService
	emailVerificationEventPublisher amqppublisher.AmqpEventPublisher
	logger                          logger.Logger
}

var _ EmailVerificationSentEventHandler = (*emailVerificationSentEventHandler)(nil)

func NewEmailVerificationSentEventHandler(
	userVerificationSvc inbound.UserVerificationService,
	messagePublisher amqppublisher.AmqpEventPublisher,
	logger logger.Logger,
) EmailVerificationSentEventHandler {
	return &emailVerificationSentEventHandler{
		userVerificationSvc:             userVerificationSvc,
		emailVerificationEventPublisher: messagePublisher,
		logger:                          logger,
	}
}

func (h *emailVerificationSentEventHandler) Handle(ctx context.Context, event events.EmailVerificationSent) error {
	h.logger.Infof("Received event email verification sent, %v", event)

	uuid, email, name := event.UserUUID, event.Email, event.Name

	verification, err := h.userVerificationSvc.CreateEmailVerification(ctx, uuid.String(), email)
	if err != nil {
		msg := fmt.Sprintf("Failed to create email verification for user %s for email %s", uuid.String(), email)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	sendEmailEvent := events.EmailVerificationSent{
		UserUUID: uuid,
		Email:    email,
		Name:     name,
		Code:     verification.Code(),
	}

	eventBytes, err := events.EventToBytes(sendEmailEvent)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse event %v", sendEmailEvent)
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
