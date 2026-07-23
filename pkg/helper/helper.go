package helper

import (
	"math/rand"
)

func RandomItemFromSlice[T any](slice []T) T {
	n := 0
	if len(slice) > 1 {
		n = rand.Intn(len(slice))
	}
	return slice[n]
}
func RandomItemFromSliceWithIndex[T any](slice []T) (T, int) {
	n := 0
	if len(slice) > 1 {
		n = rand.Intn(len(slice))
	}
	return slice[n], n
}

func SliceContainsAll[T comparable](slice []T, values []T) bool {
	items := make(map[T]struct{})
	for _, item := range slice {
		items[item] = struct{}{}
	}

	for _, value := range values {
		if _, found := items[value]; !found {
			return false
		}
	}

	return true
}

func SlicesIntersection[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	var result []T
	for _, v := range b {
		if _, exists := set[v]; exists {
			result = append(result, v)
			delete(set, v)
		}
	}

	if len(result) == 0 {
		return make([]T, 0, 1)
	}

	return result
}
