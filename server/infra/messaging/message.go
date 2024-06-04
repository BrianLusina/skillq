package messaging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Message is the message that is sent out to a topic
type Message struct {
	// ID is the message ID
	ID string `json:"id"`

	// Topic is the topic name this message is to be delivered on
	Topic string `json:"topic"`

	// ContentType is the content type of the message
	ContentType string `json:"contentType"`

	// Timestamp is the time the message was created
	Timestamp time.Time `json:"timestamp"`

	// Payload is the content of the message to be sent
	Payload any `json:"payload"`
}

// New creates a new message
func New(topic, contentType string, payload any) *Message {
	id := uuid.New().String()
	timestamp := time.Now()

	return &Message{
		ID:          id,
		Topic:       topic,
		ContentType: contentType,
		Payload:     payload,
		Timestamp:   timestamp,
	}
}

// String returns a stringified version of the message
func (m Message) String() string {
	return fmt.Sprintf("Message(id=%s, topic=%s, contentType=%s, timestamp=%s, Payload={%v})", m.ID, m.Topic, m.ContentType, m.Timestamp, m.Payload)
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
