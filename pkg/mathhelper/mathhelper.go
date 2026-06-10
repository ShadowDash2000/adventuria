package mathhelper

import "golang.org/x/exp/constraints"

// Mod normalized mod (0..n-1)
func Mod(a, m int) int {
	return ((a % m) + m) % m
}

// FloorDiv floor division
func FloorDiv(a, m int) int {
	r := Mod(a, m)
	return (a - r) / m
}

func Abs[T constraints.Signed](x T) T {
	return max(x, -x)
}
