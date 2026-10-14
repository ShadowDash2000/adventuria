package energy_changed

import (
	"errors"
)

type Payload struct {
	Amount *int   `json:"amount"`
	Source Source `json:"source,omitempty"`
}

func (p Payload) clone() Payload {
	var amount *int
	if p.Amount != nil {
		amount = new(*p.Amount)
	}

	return Payload{
		Amount: amount,
		Source: p.Source,
	}
}

func (p Payload) validate() error {
	if p.Amount == nil {
		return errors.New("energy_change: amount is empty")
	}
	return nil
}
