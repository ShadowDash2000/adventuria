package buy

import (
	"encoding/json"
	"fmt"
)

type cellShopValue struct {
	PriceMultiplier float32 `json:"price_multiplier"`
}

func (b *Buy) decodeValue(value string) (*cellShopValue, error) {
	decodedValue := cellShopValue{
		PriceMultiplier: 1,
	}
	if value != "" {
		if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cell value: %w", err)
		}
	}
	return &decodedValue, nil
}
