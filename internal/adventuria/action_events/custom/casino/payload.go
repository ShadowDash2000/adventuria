package casino

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

type Payload struct {
	ItemIds         []string `json:"item_ids"`
	PriceMultiplier float64  `json:"price_multiplier"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("casino: failed to unmarshal payload: %w", err)
	}

	err = payload.validate()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p Payload) clone() Payload {
	return Payload{
		ItemIds:         slices.Clone(p.ItemIds),
		PriceMultiplier: p.PriceMultiplier,
	}
}

func (p Payload) validate() error {
	if len(p.ItemIds) == 0 {
		return errors.New("casino: item ids is empty")
	}
	return nil
}
