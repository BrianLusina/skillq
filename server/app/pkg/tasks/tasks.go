package tasks

import (
	sharedkernel "github.com/BrianLusina/skillq/server/domain"
)

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
