package handlers

import "context"

// EventHandler is an event/task handler that can be used by an event or task handler to handle a given event
type EventHandler[T any] interface {
	// Handle handles the given event
	Handle(context.Context, T) error
}
