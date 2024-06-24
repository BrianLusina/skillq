package entity

import (
	"fmt"

	"github.com/BrianLusina/skillq/server/domain/id"
)

type (
	// EntityID is a composite of IDs used for an entity
	EntityID struct {
		// ID is a unique ID of an entity
		entityUUID id.UUID

		// KeyID is a unique Key ID for an entity used in sorting
		keyID id.KeyID

		// Identifier used to uniquely identifier an entity
		xid id.XID
	}

	EntityIDParams struct {
		// UUID is a unique UUID of an entity
		UUID id.UUID

		// KeyID is a unique Key ID for an entity used in sorting
		KeyID id.KeyID

		// Identifier used to uniquely identifier an entity
		XID id.XID
	}
)

// NewEntityID creates a new unique entity ID from the provided params
func NewEntityID(params EntityIDParams) EntityID {
	eid := params.UUID
	xid := params.XID
	keyId := params.KeyID

	return EntityID{
		entityUUID: eid,
		keyID:      keyId,
		xid:        xid,
	}
}

// UUID returns the UUID
func (eid EntityID) UUID() id.UUID {
	return eid.entityUUID
}

// XID returns the XID
func (eid EntityID) XID() id.XID {
	return eid.xid
}

// KeyID returns the Key ID
func (eid EntityID) KeyID() id.KeyID {
	return eid.keyID
}

// String returns a string representation of the ID
func (eid EntityID) String() string {
	return fmt.Sprintf("%s-%s-%s", &eid.xid, eid.keyID, eid.entityUUID)
}
