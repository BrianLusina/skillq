package events

import (
	sharedkernel "github.com/BrianLusina/skillq/server/domain"
	"github.com/BrianLusina/skillq/server/domain/id"
)

// EmailVerificationStarted is an event that triggers the start of email verification
type EmailVerificationStarted struct {
	sharedkernel.DomainEvent
	UserUUID id.UUID `json:"userId"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
}

func (e *EmailVerificationStarted) Identity() string {
	return string(EmailVerificationStartedName)
}

// EmailVerificationSent is an event that is triggered to signal that an email verification has been sent
type EmailVerificationSent struct {
	sharedkernel.DomainEvent
	UserUUID id.UUID `json:"userId"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
}

func (e *EmailVerificationSent) Identity() string {
	return string(EmailVerificationSentName)
}
