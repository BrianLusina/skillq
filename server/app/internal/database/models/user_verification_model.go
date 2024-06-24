package models

import (
	"fmt"
)

// UserVerificationModel represents the model of a user verification as stored in a database
type UserVerificationModel struct {
	BaseModel  BaseModel `bson:"inline"`
	Code       string    `bson:"code"`
	UserId     string    `bson:"user_id"`
	IsVerified bool      `bson:"is_verified"`
}

func (u *UserVerificationModel) String() string {
	return fmt.Sprintf("UserVerificationModel(base=%s, code=%s, userId=%s, isVerified=%v)",
		u.BaseModel.String(), u.Code, u.UserId, u.IsVerified)
}
