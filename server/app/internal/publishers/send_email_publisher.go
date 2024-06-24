package publishers

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/internal/domain/ports/outbound/publishers"
	"github.com/BrianLusina/skillq/server/app/pkg/tasks"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type sendEmailTaskPublisherAdapter struct {
	pub amqppublisher.AmqpEventPublisher
}

// NewSendEmailPublisher creates a new send email publisher
func NewSendEmailTaskPublisher(pub amqppublisher.AmqpEventPublisher) publishers.TaskPublisher[tasks.SendEmailVerification] {
	return &sendEmailTaskPublisherAdapter{
		pub: pub,
	}
}

// Publish implements publishers.SendEmailEventPublisherPort.
func (s *sendEmailTaskPublisherAdapter) Publish(ctx context.Context, message tasks.SendEmailVerification) error {
	sendEmailVerificationMessage := messaging.New(
		messaging.MessageParams{
			Topic:       message.Identity(),
			ContentType: "text/plain",
			Payload:     message,
		},
	)
	return s.pub.Publish(ctx, sendEmailVerificationMessage)
}

// Configure implements publishers.SendEmailEventPublisherPort.
func (s *sendEmailTaskPublisherAdapter) Configure(opts ...amqppublisher.Option) {
	s.pub.Configure(opts...)
}
