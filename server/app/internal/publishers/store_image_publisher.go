package publishers

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type storeImageTaskPublisherAdapter struct {
	pub amqppublisher.AmqpEventPublisher
}

// NewStoreImageTaskPublisher creates a new store image event publisher
func NewStoreImageTaskPublisher(pub amqppublisher.AmqpEventPublisher) publishers.TaskPublisher[tasks.StoreUserImage] {
	return &storeImageTaskPublisherAdapter{
		pub: pub,
	}
}

// Publish implements publishers.StoreImageEventPublisherPort.
func (s *storeImageTaskPublisherAdapter) Publish(ctx context.Context, message tasks.StoreUserImage) error {
	storeImageMessage := messaging.New(
		messaging.MessageParams{
			Topic:       message.Identity(),
			ContentType: "text/plain",
			Payload:     message,
		},
	)
	return s.pub.Publish(ctx, storeImageMessage)
}

// Configure implements publishers.StoreImageEventPublisherPort.
func (s *storeImageTaskPublisherAdapter) Configure(opts ...amqppublisher.Option) {
	s.pub.Configure(opts...)
}
