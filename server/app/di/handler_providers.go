package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/internal/handlers/taskhandlers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/clients/email"
	"github.com/BrianLusina/skillq/server/infra/logger"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage"
)

func ProvideStoreImageTaskHandler(
	storageClient storage.StorageClient,
	userRepo repositories.UserRepoPort,
	publisher amqppublisher.AmqpEventPublisher,
) handlers.EventHandler[tasks.StoreUserImage] {
	log := logger.New()
	storeImageTaskHandler := taskhandlers.NewStoreImageTaskHandler(storageClient, userRepo, publisher, log)

	return storeImageTaskHandler
}

func ProvideSendEmailVerificationTaskHandler(
	emailClient email.EmailClient,
	userVerificationSvc inbound.UserVerificationService,
	userRepo repositories.UserRepoPort,
) handlers.EventHandler[tasks.SendEmailVerification] {
	log := logger.New()
	sendEmailVerificationTaskHandler := taskhandlers.NewSendEmailVerificationTaskHandler(emailClient, userVerificationSvc, userRepo, log)
	return sendEmailVerificationTaskHandler
}
