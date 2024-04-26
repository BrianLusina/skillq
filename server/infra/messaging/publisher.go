package messaging

import "context"

// Publisher handles defines the methods used to handle publication of messages to a topic on a broker
type Publisher interface {
	// Publish publishes a message to a given topic
	Publish(ctx context.Context, queue string, body []byte) error

	// CloseChan closes connection to a broker
	CloseChan()
}