package events

import "github.com/BrianLusina/skillq/server/domain/id"

type UserEmailVerificationStarted struct {
	UserUUID id.UUID
	Email    string `json:"email"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}
