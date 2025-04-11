package helper

import "math/rand"

func RandomItemFromSlice[T any](slice []T) T {
	n := 0
	if len(slice) > 1 {
		n = rand.Intn(len(slice))
	}
	return slice[n]
}

func FilterByField[T any, K comparable](items []T, excludeKeys []K, keyFunc func(T) K) []T {
	excludeMap := make(map[K]struct{}, len(excludeKeys))
	for _, key := range excludeKeys {
		excludeMap[key] = struct{}{}
	}

	var filtered []T
	for _, item := range items {
		if _, found := excludeMap[keyFunc(item)]; !found {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
