package id

import (
	"log"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

// NewUUID generates a new ID
func NewUUID() UUID {
	return UUID(uuid.New())
}

// UUIDToString returns the string representation of the ID
func UUIDToString(id UUID) string {
	return uuid.UUID(id).String()
}

// UUIDToBytes converts an ID into a slice of Bytes
func UUIDToBytes(id UUID) ([]byte, error) {
	v, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// StringToIdentifier parses a string into an ID
func StringToUUID(idString string) (UUID, error) {
	value, err := uuid.Parse(idString)
	if err != nil {
		log.Println("failed to parse id: ", err.Error())
		return UUID{}, err
	}

	return UUID(value), err
}

// BytesToUUID parses a slice of bytes to an ID
func BytesToUUID(idBytes []byte) (UUID, error) {
	value, err := uuid.FromBytes(idBytes)
	if err != nil {
		return UUID{}, err
	}

	return UUID(value), nil
}
