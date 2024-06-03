package taskhandlers

import (
	"context"
	"fmt"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/internal/templates"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/domain/id"
	"github.com/BrianLusina/skillq/server/infra/clients/email"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/pkg/errors"
)

type sendEmailVerificationTaskHandler struct {
	emailClient         email.EmailClient
	userVerificationSvc inbound.UserVerificationService
	userRepo            repositories.UserRepoPort
	logger              logger.Logger
}

var _ handlers.EventHandler[tasks.SendEmailVerification] = (*sendEmailVerificationTaskHandler)(nil)

func NewSendEmailVerificationTaskHandler(
	emailClient email.EmailClient,
	userVerificationSvc inbound.UserVerificationService,
	userRepo repositories.UserRepoPort,
	logger logger.Logger,
) handlers.EventHandler[tasks.SendEmailVerification] {
	return &sendEmailVerificationTaskHandler{
		emailClient:         emailClient,
		userVerificationSvc: userVerificationSvc,
		userRepo:            userRepo,
		logger:              logger,
	}
}

func (h *sendEmailVerificationTaskHandler) Handle(ctx context.Context, task *tasks.SendEmailVerification) error {
	h.logger.Infof("Received task send email verification, %v", task)

	userID, email, name := task.UserUUID, task.Email, task.Name

	userUUID, err := id.StringToUUID(userID)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse user ID %s", userID)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	if _, err := h.userRepo.GetUserByUUID(ctx, userUUID); err != nil {
		msg := fmt.Sprintf("Failed to retrieve user %s", userID)
		h.logger.Errorf(msg)
		return errors.Wrapf(err, msg)
	}

	verification, err := h.userVerificationSvc.CreateEmailVerification(ctx, userID, email)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create verification for user %s with error %v", userID, err)
		h.logger.Error(errMsg)
		return errors.Wrapf(err, "failed to create verification for user %s with error %v", userID, err)
	}

	emailTemplate := templates.BuildEmailVerification(email, name, verification.Code())
	err = h.emailClient.Send(email, emailTemplate)
	if err != nil {
		errMsg := fmt.Sprintf("failed to send email verification for user %s with error %v", userID, err)
		h.logger.Error(errMsg)
		return errors.Wrapf(err, "failed to send email verification for user %s with error %v", userID, err)
	}

	return nil
}
