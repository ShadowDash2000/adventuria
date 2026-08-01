package model

type ActionState struct {
	Activities     *ActionActivitiesState
	ActivityFilter *ActivityFilter
	Items          *ActionItemsState
	Shop           *ActionShopState
	ShopFilter     *ActionShopFilterState
	Dealer         *ActionDealerState
}

func (a ActionState) Clone() ActionState {
	return ActionState{
		Activities:     a.Activities.Clone(),
		ActivityFilter: new(a.ActivityFilter.Clone()),
		Items:          a.Items.Clone(),
		Shop:           a.Shop.Clone(),
		ShopFilter:     a.ShopFilter.Clone(),
		Dealer:         a.Dealer.Clone(),
	}
}
