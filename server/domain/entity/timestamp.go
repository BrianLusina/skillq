package entity

import (
	"fmt"
	"time"
)

// EntityTimestamp are the timestamps for when an entity was created, updated and/or deleted
type (
	EntityTimestamp struct {
		// createdAt is when an entity was created
		createdAt time.Time

		// updatedAt is when an entity was last updated
		updatedAt time.Time

		// deletedAt is when an entity was deleted/removed from the system
		deletedAt *time.Time
	}

	// EntityTimestampParams with fields to create a new entity timestamp
	EntityTimestampParams struct {
		// CreatedAt is when an entity was created
		CreatedAt time.Time

		// UpdatedAt is when an entity was last updated
		UpdatedAt time.Time

		// DeletedAt is when an entity was deleted/removed from the system
		DeletedAt *time.Time
	}
)

// NewEntityTimestamp creates an entity timestamp
func NewEntityTimestamp(params EntityTimestampParams) EntityTimestamp {
	return EntityTimestamp{
		createdAt: params.CreatedAt,
		updatedAt: params.UpdatedAt,
		deletedAt: params.DeletedAt,
	}
}

// CreatedAt returns the timestamp this entity was created
func (et *EntityTimestamp) CreatedAt() time.Time {
	return et.createdAt
}

// UpdatedAt returns the timestamp this entity was updated
func (et *EntityTimestamp) UpdatedAt() time.Time {
	return et.updatedAt
}

// Deleted returns the timestamp this entity was deleted
func (et *EntityTimestamp) DeletedAt() *time.Time {
	return et.deletedAt
}

// ParseEntityTimestamps creates an EntityTimestamp with validation handled, if validation fails, an error is returned
func ParseEntityTimestamps(createdAt, updatedAt time.Time, deletedAt *time.Time) (EntityTimestamp, error) {
	// validate timestamps

	if createdAt.After(updatedAt) {
		// invalid, created at time should not be after updated at timestamp
		return EntityTimestamp{}, fmt.Errorf("created time %s can not be after updated time %s", createdAt, updatedAt)
	}

	// when there is a deleted at timestamp
	if deletedAt != nil {
		// it can not be before created nor updated at timestamps
		if deletedAt.Before(createdAt) || deletedAt.Before(updatedAt) {
			return EntityTimestamp{}, fmt.Errorf("created time %s can not be after updated time %s", createdAt, updatedAt)
		}
	}

	return EntityTimestamp{
		createdAt: createdAt.UTC().Round(time.Microsecond),
		updatedAt: updatedAt.UTC().Round(time.Microsecond),
		deletedAt: deletedAt,
	}, nil
}
