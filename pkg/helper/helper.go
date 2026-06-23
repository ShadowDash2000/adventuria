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
