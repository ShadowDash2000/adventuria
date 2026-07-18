package model

type ActionDataList struct {
	Activities ActivitiesData
	Items      ItemsData
}

type ActivitiesData struct {
	Ids []string
}

type ItemsData struct {
	Ids             []string
	PriceMultiplier float64
}
