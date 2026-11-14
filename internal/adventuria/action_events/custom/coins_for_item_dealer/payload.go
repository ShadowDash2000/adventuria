package coins_for_item_dealer

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	Coins       int    `json:"coins"`
	ItemId      string `json:"item_id"`
	Description string `json:"description"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("coins_for_item_dealer: failed to unmarshal payload: %w", err)
	}

	err = payload.validate()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p Payload) clone() Payload {
	return Payload{
		Coins:       p.Coins,
		ItemId:      p.ItemId,
		Description: p.Description,
	}
}

func (p Payload) validate() error {
	return nil
}
