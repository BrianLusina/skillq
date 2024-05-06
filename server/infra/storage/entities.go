package storage

import "time"

// StorageItem is the item that will be stored in blob storage
type StorageItem struct {
	// Name is the name of the document
	Name string

	// Content is the content of the storage item which is a base64 string
	Content string

	// ContentType is the type of content of this item
	ContentType string

	// Bucket is where to store the document
	Bucket string

	// Metadata is optional additional key value pair data
	Metadata map[string]string
}

// Document structure represents a document representation
type Document struct {
	// MimeType is the type of the document, application/zip, text/plain, application/pdf, image/png
	MimeType string

	// FileExtension is the document's file extension
	FileExtension string

	// Data is the data contained in the document
	Data []byte
}

type StorageDownloadItem struct {
	// Name is the name of the document
	Name string

	// Content is the content of the storage item which is a slice of bytes. The call site can use this to write to a file
	Content []byte

	// ContentType is the type of content of this item
	ContentType string

	// Bucket is where to store the document
	Bucket string

	// Metadata is optional additional key value pair data
	Metadata map[any]any

	// LastModified is the last time this item was updated/modified
	LastModified *time.Time
}
