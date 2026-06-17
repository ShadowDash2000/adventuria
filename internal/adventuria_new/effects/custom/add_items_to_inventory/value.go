package add_items_to_inventory

import (
	"strings"
)

func (a *AddItemsToInventory) decodeValue(value string) []string {
	return strings.Split(value, ";")
}
