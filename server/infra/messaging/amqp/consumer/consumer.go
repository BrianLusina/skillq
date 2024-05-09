package amqpconsumer

import (
	"context"

	"github.com/BrianLusina/skillq/server/infra/messaging"
)

// Worker function that handles consumption of messages from a given channel. This is made generic in order to work with different types
// of messages
type Worker[T any] func(ctx context.Context, message T)

// AmqpEventConsumer defines a consumer that handles consumption of messages from a Broker
type AmqpEventConsumer interface {
	messaging.EventConsumer

	// Configures an AMQP Event Consumer
	Configure(...Option) AmqpEventConsumer
}
