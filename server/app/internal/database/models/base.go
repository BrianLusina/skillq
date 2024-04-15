package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel represents the base model
type BaseModel struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty"`
	UUID      string             `bson:"uuid,omitempty"`
	XID       string             `bson:"xid,omitempty"`
	KeyID     string             `bson:"keyId,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}
