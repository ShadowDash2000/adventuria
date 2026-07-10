package teleport

import (
	"encoding/json"
	"fmt"
)

type cellTeleportValue struct {
	CellId string `json:"cell_id"`
}

func (c *CellTeleport) decodeValue(value string) (*cellTeleportValue, error) {
	var decodedValue cellTeleportValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cell value: %w", err)
	}
	return &decodedValue, nil
}
