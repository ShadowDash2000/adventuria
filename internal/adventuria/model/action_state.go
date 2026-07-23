package model

type ActionState struct {
	Activities ActionActivitiesState
	Items      ActionItemsState
	Shop       ActionShopState
	Dealer     *ActionDealerState
}

type ActionActivitiesState struct {
	Ids []string
}

type ActionItemsState struct {
	Ids []string
}

type ActionShopState struct {
	Ids             []string
	PriceMultiplier float64
}
