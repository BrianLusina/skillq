package messaging

import "context"

// EventConsumer defines a consumer that handles consumption of messages from a Broker
type EventConsumer interface {
	// Consumes a message from a given queue. This is mostly a blocking operation
	Consume(ctx context.Context, queue string) error

	// StartConsumer starts a new consumer worker. Used for async workflows
	StartConsumer(fn func(ctx context.Context, message <-chan any)) error

	// AddHandler adds a handler that will handle consumption of messages from a queue
	AddHandler(ctx context.Context, task string, handler func(payload []byte) error)
}
