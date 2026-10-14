package change_energy

import "errors"

type Payload struct {
	ProgressId string `json:"progress_id"`
	Amount     *int   `json:"amount"`
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
		return errors.New("change_energy: progress id is empty")
	}
	if p.Amount == nil {
		return errors.New("change_energy: amount is empty")
	}
	return nil
}
