package refresh_shop

import (
	"encoding/json"
	"fmt"
)

type cellShopRefreshValue struct {
	RefreshPrice int `json:"refresh_price"`
}

const defaultCellShopRefreshPrice = 10

func (r *RefreshShop) decodeValue(value string) (*cellShopRefreshValue, error) {
	decodedValue := cellShopRefreshValue{
		RefreshPrice: defaultCellShopRefreshPrice,
	}
	if value != "" {
		if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cell value: %w", err)
		}
	}
	return &decodedValue, nil
}
