package item_received

import (
	"adventuria/internal/adventuria/model"
	"encoding/json"
	"fmt"
)

type Payload struct {
	ItemId   string         `json:"item_id"`
	ItemName string         `json:"item_name"`
	ItemType model.ItemType `json:"item_type"`
}

func (i *ItemReceived) decodePayload(payload string) (*Payload, error) {
	decodedPayload := Payload{}
	if err := json.Unmarshal([]byte(payload), &decodedPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player_event payload: %w", err)
	}
	return &decodedPayload, nil
}
