package messaging

import "context"

// Worker function that handles consumption of messages from a given channel. This is made generic in order to work with different types
// of messages
type Worker[T any] func(ctx context.Context, message T)

// EventConsumer defines a consumer that handles consumption of messages from a Broker
type EventConsumer interface {
	// Consumes a message from a given queue. This is mostly a blocking operation
	Consume(ctx context.Context, queue string) error

	// StartConsumer starts a new consumer worker. Used for async workflows
	StartConsumer(fn Worker[any]) error

	// AddHandler adds a handler that will handle consumption of messages from a queue
	AddHandler(ctx context.Context, task string, handler func(payload []byte) error)
}
