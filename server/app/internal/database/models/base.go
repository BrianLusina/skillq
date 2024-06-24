package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel represents the base model
type BaseModel struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty"`
	UUID      string             `bson:"uuid,omitempty"`
	XID       string             `bson:"xid,omitempty"`
	KeyID     string             `bson:"keyId,omitempty"`
	Metadata  map[string]any     `bson:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty"`
}

func (b *BaseModel) String() string {
	return fmt.Sprintf("BaseModel(objectId=%s, uuid=%s, xid=%s, keyId=%s, metadata=%v, createdAt=%s, updatedAt=%s, deletedAt=%s)",
		b.ObjectID.String(), b.UUID, b.XID, b.KeyID, b.Metadata, b.CreatedAt, b.UpdatedAt, b.DeletedAt)
}
