package model

type ItemType string

var (
	ItemTypeBuff    ItemType = "buff"
	ItemTypeDebuff  ItemType = "debuff"
	ItemTypeNeutral ItemType = "neutral"
	ItemTypeDev     ItemType = "dev"
)

var itemTypes = map[ItemType]struct{}{
	ItemTypeBuff:    {},
	ItemTypeDebuff:  {},
	ItemTypeNeutral: {},
	ItemTypeDev:     {},
}

func ItemTypeExists(t ItemType) bool {
	_, ok := itemTypes[t]
	return ok
}
