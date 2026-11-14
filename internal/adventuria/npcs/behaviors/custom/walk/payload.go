package walk

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	MinDistance int `json:"min_distance"`
	MaxDistance int `json:"max_distance"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("walk: failed to unmarshal payload: %w", err)
	}

	err = payload.validate()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p Payload) clone() Payload {
	return Payload{
		MinDistance: p.MinDistance,
		MaxDistance: p.MaxDistance,
	}
}

func (p Payload) validate() error {
	if p.MinDistance <= 0 {
		return errors.New("walk: min_distance must be greater than zero")
	}
	if p.MaxDistance <= p.MinDistance {
		return errors.New("walk: max_distance must be greater than min_distance")
	}
	return nil
}
