package tasks

import (
	"fmt"
	"time"

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

type (
	// SendEmailVerification is a task that is triggered to signal that an email verification is to be sent
	SendEmailVerification struct {
		sharedkernel.DomainEvent
		Timestamp time.Time `json:"createdAt"`
		UserUUID  string    `json:"userId"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Code      string    `json:"code"`
	}

	SendEmailVerificationParams struct {
		UserUUID string
		Email    string
		Name     string
		Code     string
	}
)

// NewSendEmailVerificationEvent creates a new SendEmailVerification Event
func NewSendEmailVerificationEvent(params SendEmailVerificationParams) SendEmailVerification {
	timestamp := time.Now()
	return SendEmailVerification{
		Timestamp: timestamp,
		UserUUID:  params.UserUUID,
		Email:     params.Email,
		Name:      params.Name,
	}
}

func (sev *SendEmailVerification) CreatedAt() time.Time {
	return sev.Timestamp
}

func (sev *SendEmailVerification) Identity() string {
	return string(SendEmailVerificationName)
}

func (sev *SendEmailVerification) String() string {
	return fmt.Sprintf("SendEmailVerification(userUUID=%s, email=%s, name=%s)", sev.UserUUID, sev.Email, sev.Name)
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

func (st *StoreUserImageTask) Identity() string {
	return string(StoreUserImageTaskName)
}

func (st *StoreUserImageTask) String() string {
	return fmt.Sprintf("StoreUserImage(userUUID=%s, contentType=%s, name=%s, bucket=%s)", st.UserUUID, st.ContentType, st.Name, st.Bucket)
}
