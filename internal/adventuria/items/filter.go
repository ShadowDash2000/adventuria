package items

import "adventuria/internal/adventuria/model"

type Filter struct {
	ItemType         model.ItemType
	Ids              []string
	IsRollable       *bool
	PriceGreaterThan *int
}
