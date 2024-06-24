package utils

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// MessageDataToBytes marshalls an event/task message to a byte slice
func MessageDataToBytes(message any) ([]byte, error) {
	eventBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message to bytes")
	}
	return eventBytes, nil
}
