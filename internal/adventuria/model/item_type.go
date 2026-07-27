package model

import "fmt"

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

func NewItemType(value string) (ItemType, error) {
	it := ItemType(value)
	if _, ok := itemTypes[it]; !ok {
		return "", fmt.Errorf("unknown item type: %s", value)
	}
	return it, nil
}
