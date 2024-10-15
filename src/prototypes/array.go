package prototypes

func Includes[T any](arr []T, callback func(T) bool) bool {
	for _, v := range arr {
		if callback(v) {
			return true
		}
	}
	return false
}

func Map[T any](arr []T, callback func(T) T) []T {
	result := make([]T, len(arr))
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
