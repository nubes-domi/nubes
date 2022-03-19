package utils

func AnyMap[K string, V any](m map[K]V, f func(K, V) bool) bool {
	for k, v := range m {
		if f(k, v) {
			return true
		}
	}
	return false
}

func CollectMap[K comparable, V any, T any](m map[K]V, f func(V) T) []T {
	result := make([]T, len(m))

	i := 0
	for _, v := range m {
		result[i] = f(v)
		i++
	}

	return result
}

func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, len(m))

	i := 0
	for k := range m {
		result[i] = k
		i++
	}

	return result
}

func SelectMap[K comparable, V any](s map[K]V, f func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for key, value := range s {
		if f(key, value) {
			result[key] = value
		}
	}

	return result
}
