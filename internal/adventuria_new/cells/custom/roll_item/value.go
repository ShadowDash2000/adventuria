package roll_item

import (
	"encoding/json"
	"fmt"
)

type cellRollItemValue struct {
	ItemType string `json:"items_type"`
}

func (c *CellRollItem) decodeValue(value string) (*cellRollItemValue, error) {
	var decodedValue cellRollItemValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cell value: %w", err)
	}
	return &decodedValue, nil
}
