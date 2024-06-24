package amqppublisher

import (
	"github.com/BrianLusina/skillq/server/infra/messaging"
)

// AmqpEventPublisher defines the methods used to handle publication of messages/events to a topic on an AMQP broker
type AmqpEventPublisher interface {
	// inheriting the common methods for an event publisher
	messaging.EventPublisher

	// Configures an AMQP Event Publisher
	Configure(...Option) AmqpEventPublisher
}
