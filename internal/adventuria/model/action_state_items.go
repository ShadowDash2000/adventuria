package model

import "slices"

type ActionItemsState struct {
	Ids []string
}

func (a *ActionItemsState) Clone() *ActionItemsState {
	if a == nil {
		return nil
	}

	return &ActionItemsState{
		Ids: slices.Clone(a.Ids),
	}
}
