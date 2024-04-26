package storage

import "context"

// StorageClient defines all the capabilities of the storage client used to handle storage & retrieval of blob data
type StorageClient interface {
	// Upload uploads a new storage item and returns the URL to the stored item
	Upload(context.Context, StorageItem) (string, error)
}
