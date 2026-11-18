package coins_for_item

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	Coins  int    `json:"coins"`
	ItemId string `json:"item_id"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("coins_for_item: failed to unmarshal payload: %w", err)
	}

	err = payload.validate()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p Payload) Marshal() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("coins_for_item: failed to marshal payload: %w", err)
	}

	return b, nil
}

func (p Payload) clone() Payload {
	return Payload{
		Coins:  p.Coins,
		ItemId: p.ItemId,
	}
}

func (p Payload) validate() error {
	if p.ItemId == "" {
		return errors.New("coins_for_item: item_id is empty")
	}
	return nil
}
