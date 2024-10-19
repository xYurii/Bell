package prototypes

import "sort"

func Includes[T any](arr []T, callback func(T) bool) bool {
	for _, v := range arr {
		if callback(v) {
			return true
		}
	}
	return false
}

func FindIndex[T any](arr []T, callback func(T) bool) int {
	for i, v := range arr {
		if callback(v) {
			return i
		}
	}
	return -1
}

func Find[T any](arr []T, callback func(T) bool) T {
	for _, v := range arr {
		if callback(v) {
			return v
		}
	}
	var zero T
	return zero
}

func Map[T any, R any](arr []T, callback func(T) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = callback(v)
	}
	return result
}

func Filter[T any](arr []T, callback func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range arr {
		if callback(v) {
			result = append(result, v)
		}
	}
	return result
}

func Push[T any](arr []T, value T) []T {
	return append(arr, value)
}

func Splice[T any](arr []T, start int, deleteCount int, items ...T) []T {
	if start < 0 || start > len(arr) {
		return arr
	}

	if deleteCount < 0 || start+deleteCount > len(arr) {
		deleteCount = len(arr) - start
	}

	return append(arr[:start], append(items, arr[start+deleteCount:]...)...)
}

func SortSlice[T any](slice []T, less func(a, b T) bool, desc bool) {
	if desc {
		sort.Slice(slice, func(i, j int) bool {
			return less(slice[j], slice[i])
		})
	} else {
		sort.Slice(slice, func(i, j int) bool {
			return less(slice[i], slice[j])
		})
	}
}
