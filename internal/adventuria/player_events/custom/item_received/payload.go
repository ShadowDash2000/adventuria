package item_received

import (
	"adventuria/internal/adventuria/model"
	"errors"
)

type Payload struct {
	ItemId   string         `json:"item_id"`
	ItemName string         `json:"item_name"`
	ItemType model.ItemType `json:"item_type"`
	Source   Source         `json:"source,omitempty"`
}

func (p Payload) validate() error {
	if p.ItemId == "" {
		return errors.New("item_received: item id is empty")
	}
	if p.ItemName == "" {
		return errors.New("item_received: item name is empty")
	}
	if p.ItemType == "" {
		return errors.New("item_received: item type is empty")
	}
	return nil
}
