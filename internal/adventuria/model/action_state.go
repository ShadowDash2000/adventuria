package model

type ActionState struct {
	Activities     *ActionActivitiesState
	ActivityFilter *ActivityFilter
	Items          *ActionItemsState
	Shop           *ActionShopState
	Dealer         *ActionDealerState
}

func (a ActionState) Clone() ActionState {
	return ActionState{
		Activities:     a.Activities.Clone(),
		ActivityFilter: new(a.ActivityFilter.Clone()),
		Items:          a.Items.Clone(),
		Shop:           a.Shop.Clone(),
		Dealer:         a.Dealer.Clone(),
	}
}

type ActionActivitiesState struct {
	Ids []string
}

func (a *ActionActivitiesState) Clone() *ActionActivitiesState {
	if a == nil {
		return nil
	}

	return &ActionActivitiesState{
		Ids: a.Ids,
	}
}

type ActionItemsState struct {
	Ids []string
}

func (a *ActionItemsState) Clone() *ActionItemsState {
	if a == nil {
		return nil
	}

	return &ActionItemsState{
		Ids: a.Ids,
	}
}

type ActionShopState struct {
	Type            ShopType
	Ids             []string
	PriceMultiplier float64
}

func (a *ActionShopState) Clone() *ActionShopState {
	if a == nil {
		return nil
	}

	return &ActionShopState{
		Type:            a.Type,
		Ids:             a.Ids,
		PriceMultiplier: a.PriceMultiplier,
	}
}
