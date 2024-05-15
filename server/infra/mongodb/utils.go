package mongodb

// mapSortOrder maps the sort order to an integer value
func mapSortOrder(sortOrder SortOrder) int {
	switch sortOrder {
	case ASC:
		return 1
	case DESC:
		return -1
	default:
		return 1
	}
}
