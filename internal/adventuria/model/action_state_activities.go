package model

import "slices"

type ActionActivitiesState struct {
	Ids []string
}

func (a *ActionActivitiesState) Clone() *ActionActivitiesState {
	if a == nil {
		return nil
	}

	return &ActionActivitiesState{
		Ids: slices.Clone(a.Ids),
	}
}
