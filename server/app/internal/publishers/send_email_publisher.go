package publishers

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type sendEmailEventPublisherAdapter struct {
	pub amqppublisher.AmqpEventPublisher
}

// NewSendEmailPublisher creates a new send email publisher
func NewSendEmailEventPublisher(pub amqppublisher.AmqpEventPublisher) publishers.SendEmailEventPublisherPort {
	return &sendEmailEventPublisherAdapter{
		pub: pub,
	}
}

// Publish implements publishers.SendEmailEventPublisherPort.
func (s *sendEmailEventPublisherAdapter) Publish(ctx context.Context, message messaging.Message) error {
	return s.pub.Publish(ctx, message)
}

// Configure implements publishers.SendEmailEventPublisherPort.
func (s *sendEmailEventPublisherAdapter) Configure(opts ...amqppublisher.Option) amqppublisher.AmqpEventPublisher {
	return s.pub.Configure(opts...)
}

// Close implements publishers.SendEmailEventPublisherPort.
func (s *sendEmailEventPublisherAdapter) Close() error {
	return s.pub.Close()
}
