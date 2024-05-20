package mongodb

type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

// FilterOptions is a structure that contains filter options for a query to return many records
type FilterOptions struct {
	// Limit is the number of records to be returned from a query. If set to 0, all records are returned
	// If set to more than what the records are in the database, all records are returned, A negative value
	// will instead be set to the absolute value and that value is used as the limit, so -100 becomes 100
	Limit int

	// Offset is used to skip a number of records, so for a value of n, the first n documents are skipped
	Offset int

	// OrderBy represents the field to order collections by. This will be used as the key in the filter
	OrderBy string

	// Sort Order is the order to use in the sort, defaults to descending
	SortOrder SortOrder

	// FieldFilter is the map to apply to a filter to retrieve fields that match the given criteria
	FieldFilter map[string]map[string]string
}

// UpdateOptions is a structure that contains update options
type UpdateOptions struct {
	// Upsert sets whether to create a record if it is not found while updating
	Upsert bool

	// FieldOptions are field options to use for updating a document
	FieldOptions map[string]any

	// SetOptions are set options to use for updating a nested document in a document
	SetOptions map[string]map[string]any

	// FilterParams is the filter parameters to use for querying a document to update
	FilterParams FilterParams
}

// FilterParams are the filter parameters used for querying specific fields in a document
type FilterParams struct {
	// Key is the name of the field of the document
	Key string

	// Value is the value of the field to query the document
	Value any
}
