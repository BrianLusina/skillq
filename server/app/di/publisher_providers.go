package di

import (
	publisherPort "github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/app/internal/publishers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

// ProvideSendEmailTaskPublisher is used to create a send email verification task publisher for dependency injection
func ProvideSendEmailTaskPublisher(pub amqppublisher.AmqpEventPublisher) publisherPort.TaskPublisher[tasks.SendEmailVerification] {
	sendEmailTaskPublisher := publishers.NewSendEmailTaskPublisher(pub)
	return sendEmailTaskPublisher
}

// ProvideStoreImageTaskPublisher creates a store user image task publisher for injection
func ProvideStoreImageTaskPublisher(pub amqppublisher.AmqpEventPublisher) publisherPort.TaskPublisher[tasks.StoreUserImage] {
	storeImageTaskPublisher := publishers.NewStoreImageTaskPublisher(pub)
	return storeImageTaskPublisher
}
