package casino

import (
	"encoding/json"
	"fmt"
)

type cellCasinoValue struct {
	ItemIds         []string `json:"item_ids"`
	PriceMultiplier float32  `json:"price_multiplier"`
}

func (c *CellCasino) decodeValue(value string) (*cellCasinoValue, error) {
	var decodedValue cellCasinoValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cell value: %w", err)
	}
	return &decodedValue, nil
}
