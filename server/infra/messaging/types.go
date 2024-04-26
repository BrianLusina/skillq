package messaging

import (
	"fmt"
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
