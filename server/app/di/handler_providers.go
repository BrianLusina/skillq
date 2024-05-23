package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/inbound"
	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/repositories"
	"github.com/BrianLusina/skillq/server/app/internal/handlers"
	"github.com/BrianLusina/skillq/server/app/internal/handlers/eventhandlers"
	"github.com/BrianLusina/skillq/server/app/internal/handlers/taskhandlers"
	"github.com/BrianLusina/skillq/server/app/pkg/events"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
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

func ProvideEmailVerificationStartedEventHandler(
	userVerificationSvc inbound.UserVerificationService,
	messagePublisher amqppublisher.AmqpEventPublisher,
) handlers.EventHandler[events.EmailVerificationStarted] {
	log := logger.New()
	emailVerificationStartedHandler := eventhandlers.NewEmailVerificationStartedEventHandler(userVerificationSvc, messagePublisher, log)
	return emailVerificationStartedHandler
}
func ProvideEmailVerificationSentEventHandler(
	userVerificationSvc inbound.UserVerificationService,
	messagePublisher amqppublisher.AmqpEventPublisher,
) handlers.EventHandler[events.EmailVerificationSent] {
	log := logger.New()
	emailVerificationSentHandler := eventhandlers.NewEmailVerificationSentEventHandler(userVerificationSvc, messagePublisher, log)
	return emailVerificationSentHandler
}
