package util

func Convert[T any, P any](items []T, converterF func(T) P) []P {
	values := make([]P, 0, len(items))

	for _, item := range items {
		values = append(values, converterF(item))
	}

	return values
}

func Contains[T comparable](elems []T, elem T) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func Keys[T any, K comparable](m map[K]T) []K {
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
