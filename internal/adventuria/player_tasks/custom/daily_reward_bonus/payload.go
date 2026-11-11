package daily_reward_bonus

import (
	"errors"
)

type Payload struct {
	Amount *int `json:"amount"`
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
