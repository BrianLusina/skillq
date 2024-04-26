package storage

// StorageItem is the item that will be stored in blob storage
type StorageItem struct {
	// Type is the type of the document to store
	Type string
	// Name is the name of the document
	Name string
	// Content is the content of the storage item
	Content string
	// Bucket is where to store the document
	Bucket string
	// Metadata is optional additional key value pair data
	Metadata map[string]any
}
