package storage

// PolicyType is the type of policy to apply to a storage item
type PolicyType string

const (
	// PolicyTypeReadOnly is the read only policy type that allows only reading the storage item
	PolicyTypeReadOnly PolicyType = "READ_ONLY"

	// PolicyTypeWriteOnly is the write only policy type that allows writing to the storage item
	PolicyTypeWriteOnly PolicyType = "WRITE_ONLY"

	// PolicyTypeReadAndWrite is the read and write policy type that allows reading and writing to the storage item
	PolicyTypeReadAndWrite PolicyType = "READ_AND_WRITE"
)
