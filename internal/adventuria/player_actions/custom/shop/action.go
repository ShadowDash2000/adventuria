package shop

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions"
	"bytes"
	"fmt"
	"math/rand/v2"
	"slices"
)

const PlayerActionType model.PlayerActionType = "shop"

type Shop struct {
	*player_actions.PlayerActionBase
	payload Payload
}

func NewDef() player_actions.PlayerActionDef {
	return player_actions.NewPlayerAction(
		PlayerActionType,
		func(info *model.PlayerActionInfo) (model.PlayerActionEntry, error) {
			return NewInstance(info)
		},
	)
}

func New(payload Payload) (*Shop, error) {
	payload = payload.clone()

	if payload.RefreshPrice == 0 {
		payload.RefreshPrice = defaultRefreshPrice
	}

	if payload.MaxItems == 0 {
		payload.MaxItems = defaultMaxItems
	}

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	return &Shop{
		PlayerActionBase: player_actions.NewPlayerActionBase(PlayerActionType, encoded),
		payload:          payload,
	}, nil
}

func NewFromJSON(data []byte) (*Shop, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &Shop{
		PlayerActionBase: player_actions.NewPlayerActionBase(PlayerActionType, bytes.Clone(data)),
		payload:          payload,
	}, nil
}

func NewInstance(info *model.PlayerActionInfo) (model.PlayerActionEntry, error) {
	action, err := NewFromJSON(info.Payload())
	if err != nil {
		return nil, err
	}

	return model.NewPlayerActionInstance(info, action), nil
}

func (s *Shop) ShopType() Type {
	return s.payload.Type
}

func (s *Shop) IDs() []string {
	return slices.Clone(s.payload.Ids)
}

func (s *Shop) RefreshPrice() int {
	return s.payload.RefreshPrice
}

func (s *Shop) ApplyPriceMultiplier(basePrice int) int {
	if s.payload.PriceMultiplier == 0 {
		return basePrice
	}

	return int(float64(basePrice) * s.payload.PriceMultiplier)
}

func (s *Shop) RemoveItem(id string) error {
	payload := s.payload

	index := slices.Index(payload.Ids, id)
	if index == -1 {
		return fmt.Errorf("shop: item with id = %s not found", id)
	}
	payload.Ids = slices.Delete(slices.Clone(payload.Ids), index, index+1)

	return s.setPayload(payload)
}

func (s *Shop) HasItem(id string) bool {
	return slices.Contains(s.payload.Ids, id)
}

func (s *Shop) Refresh(ids []string) error {
	payload := s.payload

	rndIds := make([]string, 0, payload.MaxItems)
	if len(ids) > 0 {
		for range payload.MaxItems {
			rndIds = append(rndIds, ids[rand.N(len(ids))])
		}
	}

	payload.Ids = rndIds

	return s.setPayload(payload)
}

func (s *Shop) SetItemIDs(ids []string) error {
	payload := s.payload
	payload.Ids = slices.Clone(ids)

	return s.setPayload(payload)
}

func (s *Shop) Filter() Filter {
	return s.payload.Filter.clone()
}

func (s *Shop) setPayload(payload Payload) error {
	err := payload.validate()
	if err != nil {
		return err
	}

	encoded, err := payload.Marshal()
	if err != nil {
		return err
	}

	s.payload = payload
	s.PlayerActionBase.SetPayload(encoded)

	return nil
}
