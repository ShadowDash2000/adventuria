package change_balance

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	ProgressId string `json:"progress_id"`
	Amount     int    `json:"amount"`
}

func (c *ChangeBalance) decodePayload(payload string) (*Payload, error) {
	decodedPayload := Payload{}
	if err := json.Unmarshal([]byte(payload), &decodedPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal outbox payload: %w", err)
	}
	return &decodedPayload, nil
}
