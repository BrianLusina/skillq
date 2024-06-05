package publishers

import (
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type (

	// SendEmailPublisher is a publisher for sending events/tasks/messages to a queue that handles sending of emails
	SendEmailEventPublisherPort interface {
		amqppublisher.AmqpEventPublisher
	}

	// StoreImageEventPublisherPort is a publisher for sending events/tasks/messages to a queue that handles storing of images
	StoreImageEventPublisherPort interface {
		amqppublisher.AmqpEventPublisher
	}
)
