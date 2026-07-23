package casino

import (
	"encoding/json"
	"fmt"
)

type cellEventCasinoValue struct {
	ItemIds         []string `json:"item_ids"`
	PriceMultiplier float64  `json:"price_multiplier"`
}

func (c *Casino) decodeValue(value string) (*cellEventCasinoValue, error) {
	var decodedValue cellEventCasinoValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event value: %w", err)
	}
	return &decodedValue, nil
}
