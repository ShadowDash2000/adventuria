package items

import "adventuria/internal/adventuria/model"

type Filter struct {
	ItemType         model.ItemType
	Ids              []string
	Disabled         *bool
	IsRollable       *bool
	PriceGreaterThan *int
}
