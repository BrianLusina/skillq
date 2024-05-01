package tools

import "github.com/samber/lo"

// Map maps a collection/slice of objects of type T to type R with the given function fn
func Map[T any, R any](collection []T, fn func(item T, idx int) R) []R {
	return lo.Map(collection, fn)
}

// MapWithError maps a collection/slice of objects of type T to type R with the given function fn and potentially returns an error
func MapWithError[T any, R any](collection []T, fn func(item T, idx int) (R, error)) ([]R, error) {
	result := make([]R, len(collection))

	for i, item := range collection {
		r, err := fn(item, i)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}

// Filter filters data based on a certain predicate
func Filter[T any](data []T, predicate func(T) bool) []T {
	filtered := make([]T, 0)

	for _, d := range data {
		if predicate(d) {
			filtered = append(filtered, d)
		}
	}

	return filtered
}
