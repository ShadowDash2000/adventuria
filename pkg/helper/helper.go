package helper

import "math/rand"

func RandomItemFromSlice[T any](slice []T) T {
	return slice[rand.Intn(len(slice)-1)]
}
