package coins_for_item_dealer

import (
	"encoding/json"
	"fmt"
)

type cellEventCoinsForItemDealerValue struct {
	Coins       int    `json:"coins"`
	ItemId      string `json:"item_id"`
	Description string `json:"description"`
}

func (c *CoinsForItemDealer) decodeValue(value string) (*cellEventCoinsForItemDealerValue, error) {
	var decodedValue cellEventCoinsForItemDealerValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event value: %w", err)
	}
	return &decodedValue, nil
}
