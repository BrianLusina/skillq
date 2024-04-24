package events

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// EventToBytes marshalls an event message to a byte slice
func EventToBytes(message any) ([]byte, error) {
	eventBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message to bytes")
	}
	return eventBytes, nil
}
