package messaging

import (
	"encoding/json"
	"fmt"
)

// Message is the message that is sent out to a topic
type Message struct {
	Task    string `json:"task"`
	Payload any    `json:"payload"`
}

// ConsumeMessage is the message that is consumed from a queue
type ConsumeMessage struct {
	Task    string          `json:"task"`
	Payload json.RawMessage `json:"payload"`
}

// String returns a stringified version of the message
func (m Message) String() string {
	return fmt.Sprintf("Message(task=%s, Payload={%s})", m.Task, m.Payload)
}

// String returns a stringified version of the message
func (m ConsumeMessage) String() string {
	return fmt.Sprintf("ConsumeMessage(task=%s, Payload={%s})", m.Task, m.Payload)
}
