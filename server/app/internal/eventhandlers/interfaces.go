package eventhandlers

import (
	"context"

	"github.com/BrianLusina/skillq/server/app/pkg/events"
)

// EmailVerificationStartedEventHandler is an event handler that handles the user email verification started event
type EmailVerificationStartedEventHandler interface {
	// Handle handles the given event
	Handle(context.Context, events.EmailVerificationStarted) error
}

// EmailVerificationSentEventHandler is an event handler that handles the user email verification sent event
type EmailVerificationSentEventHandler interface {
	// Handle handles the given event
	Handle(context.Context, events.EmailVerificationSent) error
}
