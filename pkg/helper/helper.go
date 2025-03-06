package helper

import "math/rand"

func RandomItemFromSlice[T any](slice []T) T {
	n := 0
	if len(slice) > 1 {
		n = rand.Intn(len(slice) - 1)
	}
	return slice[n]
}
