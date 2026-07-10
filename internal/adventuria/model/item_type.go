package model

import "fmt"

type ItemType string

var (
	ItemTypeBuff    ItemType = "buff"
	ItemTypeDebuff  ItemType = "debuff"
	ItemTypeNeutral ItemType = "neutral"
)

var itemTypes = map[ItemType]struct{}{
	ItemTypeBuff:    {},
	ItemTypeDebuff:  {},
	ItemTypeNeutral: {},
}

func NewItemType(value string) (ItemType, error) {
	it := ItemType(value)
	if _, ok := itemTypes[it]; !ok {
		return "", fmt.Errorf("unknown item type: %s", value)
	}
	return it, nil
}
