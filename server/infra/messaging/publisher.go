package messaging

import "context"

// EventPublisher handles defines the methods used to handle publication of messages/events to a topic on a broker
type EventPublisher interface {
	// Publish publishes a message to a given topic
	Publish(ctx context.Context, message Message) error

	// Close closes connection to a broker
	Close()
}
