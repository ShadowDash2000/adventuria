package daily_reward_bonus

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	Amount *int `json:"amount"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("daily_reward_bonus: failed to unmarshal payload: %w", err)
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
		return nil, fmt.Errorf("daily_reward_bonus: failed to marshal payload: %w", err)
	}

	return b, nil
}

func (p Payload) clone() Payload {
	var amount *int
	if p.Amount != nil {
		amount = new(*p.Amount)
	}

	return Payload{
		Amount: amount,
	}
}

func (p Payload) validate() error {
	if p.Amount == nil {
		return errors.New("daily_reward_bonus: amount is empty")
	}
	return nil
}
