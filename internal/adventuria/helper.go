package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/exp/constraints"
)

func GetRecordById(table, id string, expand []string) (*core.Record, error) {
	record, err := PocketBase.FindRecordById(GameCollections.Get(table), id)
	if err != nil {
		return nil, err
	}

	if expand != nil {
		errs := PocketBase.ExpandRecord(record, expand, nil)
		if errs != nil {
			for _, err := range errs {
				return nil, err
			}
		}
	}

	return record, nil
}

// normalized mod (0..n-1)
func mod(a, m int) int {
	return ((a % m) + m) % m
}

// floor division
func floorDiv(a, m int) int {
	r := mod(a, m)
	return (a - r) / m
}

func abs[T constraints.Signed](x T) T {
	return max(x, -x)
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
