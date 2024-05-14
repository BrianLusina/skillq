package common

// RequestParams are parameters that can be passed when fetching records from a database
type RequestParams struct {
	// IncludedDeleted includes "softly" deleted records in the result of the query
	IncludeDeleted bool

	// Limit is the page size
	Limit int

	// Offset is the page number or cursor
	Offset int

	// OrderOption contains ordering options for the filter
	OrderOption RequestParamsOrderOption
}

// RequestParamOptions is a functional option that allows setting values to request parameters
type RequestParamOptions func(*RequestParams)

// WithRequestLimit sets the limit of a request
func WithRequestLimit(limit int) RequestParamOptions {
	return func(rp *RequestParams) {
		rp.Limit = limit
	}
}

// WithOffset sets the offset of a pagination request
func WithOffset(offset int) RequestParamOptions {
	return func(rp *RequestParams) {
		rp.Offset = offset
	}
}

// WithIncludeDeleted sets whether to include deleted records
func WithIncludeDeleted(hasDeleted bool) RequestParamOptions {
	return func(rp *RequestParams) {
		rp.IncludeDeleted = hasDeleted
	}
}

// WithSortOrder sets how to sort records
func WithSortOrder(sortOrder SortOrder) RequestParamOptions {
	return func(rp *RequestParams) {
		rp.OrderOption.SortOrder = sortOrder
	}
}

// WithSortOrder sets how to sort records
func WithOrderBy(orderBy OrderBy) RequestParamOptions {
	return func(rp *RequestParams) {
		rp.OrderOption.OrderBy = orderBy
	}
}

// NewRequestParams creates a new request params object with reasonable defaults set
func NewRequestParams(opts ...RequestParamOptions) RequestParams {
	r := RequestParams{
		Limit:          100,
		Offset:         0,
		IncludeDeleted: false,
		OrderOption: RequestParamsOrderOption{
			OrderBy:   CREATED_AT,
			SortOrder: DESC,
		},
	}

	for _, opt := range opts {
		opt(&r)
	}

	return r
}
