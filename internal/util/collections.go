package util

import (
	"sort"
	"strings"
)

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

func Join[T any](items []T, sep string, captionF func(T) string) string {
	b := strings.Builder{}

	for i, item := range items {
		b.WriteString(captionF(item))

		if i != len(items)-1 {
			b.WriteString(sep)
		}
	}

	return b.String()
}

func JoinStrings(items []string, sep string) string {
	return Join(items, sep, func(s string) string { return s })
}

func SortBy[T any](items []T, lessF func(T, T) bool) []T {
	sort.Slice(items,
		func(i, j int) bool {
			return lessF(items[i], items[j])
		},
	)
	return items
}

func FilterBy[T any](items []T, filterF func(T) bool) []T {
	result := make([]T, 0)

	for _, item := range items {
		if filterF(item) {
			result = append(result, item)
		}
	}

	return result
}
