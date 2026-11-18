package coins_for_item

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions"
	"bytes"
)

const PlayerActionType model.PlayerActionType = "coins_for_item"

type CoinsForItem struct {
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

func New(payload Payload) (*CoinsForItem, error) {
	payload = payload.clone()

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	return &CoinsForItem{
		PlayerActionBase: player_actions.NewPlayerActionBase(PlayerActionType, encoded),
		payload:          payload,
	}, nil
}

func NewFromJSON(data []byte) (*CoinsForItem, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &CoinsForItem{
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

func (c *CoinsForItem) Coins() int {
	return c.payload.Coins
}

func (c *CoinsForItem) ItemID() string {
	return c.payload.ItemId
}
