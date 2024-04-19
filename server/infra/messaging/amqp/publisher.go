package amqp

import "context"

// AmqpPublisher handles defines the methods used to handle publication of messages to a topic on a broker
type AmqpPublisher interface {
	// Publish publishes a message to a given topic
	Publish(ctx context.Context, queue string, task string, body any) error

	// CloseChan closes connection to an AMQP broker
	CloseChan()
}
