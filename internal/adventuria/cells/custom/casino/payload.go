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

func NewPayloadFromJSON(payload []byte) (Payload, error) {
	var decoded Payload
	err := json.Unmarshal(payload, &decoded)
	if err != nil {
		return Payload{}, fmt.Errorf("casino: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return Payload{}, err
	}

	return decoded, nil
}

func (p Payload) clone() Payload {
	return Payload{
		ItemIds:         slices.Clone(p.ItemIds),
		PriceMultiplier: p.PriceMultiplier,
	}
}

func (p Payload) validate() error {
	if len(p.ItemIds) == 0 {
		return errors.New("casino: item_ids is empty")
	}
	return nil
}
