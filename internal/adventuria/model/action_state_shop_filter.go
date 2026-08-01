package model

import "slices"

type ActionShopFilterState struct {
	ItemType ItemType
	Ids      []string
}

func (a *ActionShopFilterState) Clone() *ActionShopFilterState {
	if a == nil {
		return nil
	}

	return &ActionShopFilterState{
		ItemType: a.ItemType,
		Ids:      slices.Clone(a.Ids),
	}
}
