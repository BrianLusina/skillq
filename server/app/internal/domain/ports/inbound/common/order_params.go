package common

// SortOrder to use when sorting records
type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

// OrderBy is the field to order collections by
type OrderBy string

const (
	CREATED_AT OrderBy = "created_at"
	UPDATED_AT OrderBy = "updated_at"
	DELETED_AT OrderBy = "deleted_at"
)

// RequestParamsOrderOption contains filtering values to use for ordering
type RequestParamsOrderOption struct {
	// OrderBy is the field to order records by
	OrderBy OrderBy

	// SortOrder is the order to sort the records by
	SortOrder SortOrder
}
