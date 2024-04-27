package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserVerificationModel represents the model of a user verification as stored in a database
type UserVerificationModel struct {
	ObjectID   primitive.ObjectID `bson:"_id,omitempty"`
	UUID       string             `bson:"uuid,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	Code       string             `bson:"code"`
	UserId     string             `bson:"user_id"`
	IsVerified bool               `bson:"is_verified"`
}
