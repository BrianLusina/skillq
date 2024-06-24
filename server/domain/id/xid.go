package id

import (
	"log"

	"github.com/rs/xid"
)

type XID = xid.ID

// NewXid generates a new ID
func NewXid() XID {
	return XID(xid.New())
}

// XidToString returns the string representation of the ID
func XidToString(id XID) string {
	return xid.ID(id).String()
}

// XidToBytes converts an ID into a slice of Bytes
func XidToBytes(id XID) []byte {
	return xid.ID(id).Bytes()
}

// StringToXid parses a string into an ID
func StringToXid(idString string) (XID, error) {
	value, err := xid.FromString(idString)
	if err != nil {
		log.Println("failed to parse id: ", err.Error())
	}

	return XID(value), err
}

// BytesToXid parses a slice of bytes to an ID
func BytesToXid(idString []byte) (XID, error) {
	value, err := xid.FromBytes(idString)
	if err != nil {
		return XID{}, err
	}

	return XID(value), nil
}
