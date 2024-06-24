package publishers

import (
	"context"

	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
)

type (
	// TaskPublisher publishes a task that is to be processed and used
	TaskPublisher[T any] interface {
		Publish(ctx context.Context, message T) error
		Configure(...amqppublisher.Option)
	}
)
