package helper

import "math/rand"

func RandomItemFromSlice[T any](slice []T) T {
	n := 0
	if len(slice) > 1 {
		n = rand.Intn(len(slice))
	}
	return slice[n]
}
