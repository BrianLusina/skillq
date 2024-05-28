package eventhandlers

import (
	"context"
	"fmt"

	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/pkg/errors"
)

type emailVerificationStartedEventHandler struct {
	userVerificationSvc             inbound.UserVerificationService
	emailVerificationEventPublisher amqppublisher.AmqpEventPublisher
	logger                          logger.Logger
}

var _ handlers.EventHandler[events.EmailVerificationStarted] = (*emailVerificationStartedEventHandler)(nil)

func NewEmailVerificationStartedEventHandler(
	userVerificationSvc inbound.UserVerificationService,
	messagePublisher amqppublisher.AmqpEventPublisher,
	logger logger.Logger,
) handlers.EventHandler[events.EmailVerificationStarted] {
	return &emailVerificationStartedEventHandler{
		userVerificationSvc:             userVerificationSvc,
		emailVerificationEventPublisher: messagePublisher,
		logger:                          logger,
	}
}

func (h *emailVerificationStartedEventHandler) Handle(ctx context.Context, event *events.EmailVerificationStarted) error {
	h.logger.Infof("Received event email verification started, %v", event)

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

	sendEmailVerificationMessage := messaging.Message{
		Topic:       sendEmailEvent.Identity(),
		ContentType: "text/plain",
		Payload:     sendEmailEvent,
	}

	if err := h.emailVerificationEventPublisher.Publish(ctx, sendEmailVerificationMessage); err != nil {
		msg := fmt.Sprintf("Failed to publish event %v", sendEmailVerificationMessage)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, "failed to publish user email sent event")
	}

	return nil
}
