package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Message is the message that is sent out to a topic
type Message struct {
	Task    string `json:"task"`
	Payload any    `json:"payload"`
}

// String returns a stringified version of the message
func (m Message) String() string {
	return fmt.Sprintf("Message(task=%s, Payload={%s})", m.Task, m.Payload)
}

// ToBytes converts a given message to bytes
func ToBytes(message Message) ([]byte, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert message %s to bytes", message)
	}

	return bytes, nil
}
