package change_balance

import (
	"encoding/json"
	"fmt"
)

type OutboxValue struct {
	ProgressId string `json:"progress_id"`
	Amount     int    `json:"amount"`
}

func (c *ChangeBalance) decodeValue(value string) (*OutboxValue, error) {
	decodedValue := OutboxValue{}
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal outbox value: %w", err)
	}
	return &decodedValue, nil
}
