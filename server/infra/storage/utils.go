package storage

import "github.com/vincent-petithory/dataurl"

// GetDocumentData parses a base64 string into a document
func GetDocumentData(base64 string) (*Document, error) {
	doc, err := dataurl.DecodeString(base64)
	if err != nil {
		return nil, err
	}

	return &Document{
		MimeType:      doc.MediaType.ContentType(),
		Data:          doc.Data,
		FileExtension: doc.Subtype,
	}, nil
}
