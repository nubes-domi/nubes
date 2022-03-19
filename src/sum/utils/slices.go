package utils

func Any[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func Collect[T any, R any](s []T, f func(T) R) []R {
	result := make([]R, len(s))

	i := 0
	for _, v := range s {
		result[i] = f(v)
		i++
	}

	return result
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Select[T any](s []T, f func(T) bool) []T {
	result := []T{}
	for _, value := range s {
		if f(value) {
			result = append(result, value)
		}
	}

	return result
}
