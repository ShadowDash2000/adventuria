package model

type ActionState struct {
	Activities     *ActionActivitiesState
	UsedItems      ActionUsedItemsState
	ActivityFilter *ActivityFilter
	Items          *ActionItemsState
	Shop           *ActionShopState
	ShopFilter     *ActionShopFilterState
	Dealer         *ActionDealerState
}

func (a ActionState) Clone() ActionState {
	return ActionState{
		Activities:     a.Activities.Clone(),
		UsedItems:      a.UsedItems.Clone(),
		ActivityFilter: a.ActivityFilter.CloneNil(),
		Items:          a.Items.Clone(),
		Shop:           a.Shop.Clone(),
		ShopFilter:     a.ShopFilter.Clone(),
		Dealer:         a.Dealer.Clone(),
	}
}
