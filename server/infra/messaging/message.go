package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Message is the message that is sent out to a topic
type Message struct {
	// Topic is the topic name this message is to be delivered on
	Topic string `json:"topic"`
	// ContentType is the content type of the message
	ContentType string `json:"contentType"`
	// Payload is the content of the message to be sent
	Payload any `json:"payload"`
}

// String returns a stringified version of the message
func (m Message) String() string {
	return fmt.Sprintf("Message(task=%s, Payload={%v})", m.Topic, m.Payload)
}

// ToBytes marshalls an event/task message to a byte slice
func (m *Message) ToBytes() ([]byte, error) {
	eventBytes, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message to bytes")
	}
	return eventBytes, nil
}

// ConsumeMessage is the message that is consumed from a queue
type ConsumeMessage struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// String returns a stringified version of the message
func (m ConsumeMessage) String() string {
	return fmt.Sprintf("ConsumeMessage(task=%s, Payload={%v})", m.Topic, m.Payload)
}
