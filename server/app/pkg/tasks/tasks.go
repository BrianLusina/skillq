package tasks

import (
	sharedkernel "github.com/BrianLusina/skillq/server/domain"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// StartEmailVerification is an event that triggers the start of email verification
type StartEmailVerification struct {
	sharedkernel.DomainEvent
	UserUUID id.UUID `json:"userId"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
}

func (e *StartEmailVerification) Identity() string {
	return string(StartEmailVerificationName)
}

// SendEmailVerification is a task that is triggered to signal that an email verification is to be sent
type SendEmailVerification struct {
	sharedkernel.DomainEvent
	UserUUID string `json:"userId"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

func (e *SendEmailVerification) Identity() string {
	return string(SendEmailVerificationName)
}

// StoreUserImageTask is a task message that triggers the storage of a user image
type StoreUserImageTask struct {
	sharedkernel.DomainEvent
	UserUUID    string `json:"userUUID"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Name        string `json:"name"`
	Bucket      string `json:"bucket"`
}

func (e *StoreUserImageTask) Identity() string {
	return string(StoreUserImageTaskName)
}
