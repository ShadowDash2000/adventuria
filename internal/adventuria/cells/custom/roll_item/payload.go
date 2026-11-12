package roll_item

import (
	"adventuria/internal/adventuria/model"
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	ItemsType string `json:"items_type"`
}

func NewPayloadFromJSON(payload []byte) (Payload, error) {
	var decoded Payload
	err := json.Unmarshal(payload, &decoded)
	if err != nil {
		return Payload{}, fmt.Errorf("roll_item: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return Payload{}, err
	}

	return decoded, nil
}

func (p Payload) clone() Payload {
	return Payload{
		ItemsType: p.ItemsType,
	}
}

func (p Payload) validate() error {
	if p.ItemsType == "" {
		return errors.New("roll_item: items_type is empty")
	}
	if !model.ItemTypeExists(model.ItemType(p.ItemsType)) {
		return errors.New("roll_item: items_type is invalid")
	}
	return nil
}
