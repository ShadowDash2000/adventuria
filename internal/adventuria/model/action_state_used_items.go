package model

import "slices"

type ActionUsedItemsState []ActionUsedItemState

func (a ActionUsedItemsState) Clone() ActionUsedItemsState {
	if a == nil {
		return nil
	}

	return slices.Clone(a)
}

type ActionUsedItemState struct {
	Id string
}
