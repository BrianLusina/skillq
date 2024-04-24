package amqp

import "context"

// Amqp defines an interface that Advanced Message Queuing Protocol clients implement
type Amqp interface {
	// Declare creates a given queue with a given message TTL(Time To Live) and also declares a DLX(Dead Letter Exchange) on the queue
	Declare(ctx context.Context, queue string, msgTTL int, dlx bool) error

	// Publish publishes a message to a given topic
	Publish(ctx context.Context, queue string, body any) error

	// Consumes a message from a given queue
	Consume(ctx context.Context, queue string) error

	// AddHandler adds a handler that will handle consumption of messages from a queue
	AddHandler(ctx context.Context, task string, handler func(payload []byte) error)
}
