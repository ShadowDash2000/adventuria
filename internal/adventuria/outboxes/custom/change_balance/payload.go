package change_balance

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	ProgressId string `json:"progress_id"`
	Amount     *int   `json:"amount"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("change_balance: failed to unmarshal payload: %w", err)
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
		return nil, fmt.Errorf("change_balance: failed to marshal payload: %w", err)
	}

	return b, nil
}

func (p Payload) clone() Payload {
	var amount *int
	if p.Amount != nil {
		amount = new(*p.Amount)
	}

	return Payload{
		ProgressId: p.ProgressId,
		Amount:     amount,
	}
}

func (p Payload) validate() error {
	if p.ProgressId == "" {
		return errors.New("change_balance: progress id is empty")
	}
	if p.Amount == nil {
		return errors.New("change_balance: amount is empty")
	}
	return nil
}
