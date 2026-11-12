package teleport

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	CellId string `json:"cell_id"`
}

func NewPayloadFromJSON(payload []byte) (Payload, error) {
	var decoded Payload
	err := json.Unmarshal(payload, &decoded)
	if err != nil {
		return Payload{}, fmt.Errorf("teleport: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return Payload{}, err
	}

	return decoded, nil
}

func (p Payload) clone() Payload {
	return Payload{
		CellId: p.CellId,
	}
}

func (p Payload) validate() error {
	if p.CellId == "" {
		return errors.New("teleport: cell_id is empty")
	}
	return nil
}
