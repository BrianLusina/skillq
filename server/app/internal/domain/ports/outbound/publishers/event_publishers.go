package publishers

import (
	"context"

	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type (
	// EventPublisher publishes events that have already happened in the system and that may be useful for other consumers to listen on
	EventPublisher[T any] interface {
		Publish(ctx context.Context, message T) error
		Configure(...amqppublisher.Option)
	}
)
