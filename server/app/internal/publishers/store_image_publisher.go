package publishers

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type storeImageEventPublisherAdapter struct {
	pub amqppublisher.AmqpEventPublisher
}

// NewStoreImageEventPublisher creates a new store image event publisher
func NewStoreImageEventPublisher(pub amqppublisher.AmqpEventPublisher) publishers.StoreImageEventPublisherPort {
	return &storeImageEventPublisherAdapter{
		pub: pub,
	}
}

// Publish implements publishers.StoreImageEventPublisherPort.
func (s *storeImageEventPublisherAdapter) Publish(ctx context.Context, message messaging.Message) error {
	return s.pub.Publish(ctx, message)
}

// Configure implements publishers.StoreImageEventPublisherPort.
func (s *storeImageEventPublisherAdapter) Configure(opts ...amqppublisher.Option) amqppublisher.AmqpEventPublisher {
	return s.pub.Configure(opts...)
}

// Close implements publishers.StoreImageEventPublisherPort.
func (s *storeImageEventPublisherAdapter) Close() error {
	return s.pub.Close()
}
