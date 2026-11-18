package shop

import (
	"adventuria/internal/adventuria/model"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

type Type string

const (
	TypeBuffet Type = "buffet"
	TypeCasino Type = "casino"
)

const (
	defaultRefreshPrice = 10
	defaultMaxItems     = 6
)

type Payload struct {
	Type            Type     `json:"type"`
	Ids             []string `json:"ids,omitempty"`
	PriceMultiplier float64  `json:"price_multiplier"`
	RefreshPrice    int      `json:"refresh_price"`
	MaxItems        int      `json:"max_items"`
	Filter          Filter   `json:"filter"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("shop: failed to unmarshal payload: %w", err)
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
		return nil, fmt.Errorf("shop: failed to marshal payload: %w", err)
	}

	return b, nil
}

func (p Payload) clone() Payload {
	return Payload{
		Type:            p.Type,
		Ids:             slices.Clone(p.Ids),
		PriceMultiplier: p.PriceMultiplier,
		RefreshPrice:    p.RefreshPrice,
		MaxItems:        p.MaxItems,
		Filter:          p.Filter.clone(),
	}
}

func (p Payload) validate() error {
	if p.Type == "" {
		return errors.New("shop: type is empty")
	}
	return p.Filter.validate()
}

type Filter struct {
	ItemType model.ItemType `json:"item_type,omitempty"`
	Ids      []string       `json:"ids,omitempty"`
}

func (f Filter) clone() Filter {
	return Filter{
		ItemType: f.ItemType,
		Ids:      slices.Clone(f.Ids),
	}
}

func (f Filter) validate() error {
	if f.ItemType == "" && f.Ids == nil {
		return errors.New("shop: filter needs at least one field")
	}
	return nil
}
