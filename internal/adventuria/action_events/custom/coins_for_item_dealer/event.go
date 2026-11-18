package coins_for_item_dealer

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/deals/coins_for_item"
	"context"
)

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

type playerActions interface {
	Save(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
	DeleteByNPC(ctx context.Context, npc *model.NPC, player *model.Player) error
}

var _ action_events.ActionEvent = (*CoinsForItemDealer)(nil)

const Type model.ActionEventType = "coins_for_item_dealer"

type CoinsForItemDealer struct {
	items         items
	payload       Payload
	playerActions playerActions
}

func NewDef(items items, playerActions playerActions) action_events.ActionEventDef {
	return action_events.NewActionEventDef(
		Type,
		func(data []byte) (action_events.ActionEvent, error) {
			return NewFromJSON(items, playerActions, data)
		},
	)
}

func NewFromJSON(items items, playerActions playerActions, data []byte) (*CoinsForItemDealer, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &CoinsForItemDealer{
		items:         items,
		payload:       payload,
		playerActions: playerActions,
	}, nil
}

func (c *CoinsForItemDealer) Init(ctx context.Context, npc *model.NPC, player *model.Player) error {
	dealAction, err := coins_for_item.New(coins_for_item.Payload{
		Coins:  c.payload.Coins,
		ItemId: c.payload.ItemId,
	})
	if err != nil {
		return err
	}

	action, err := model.NewPlayerAction(model.PlayerActionCreate{
		PlayerId:      player.ID(),
		ActionId:      player.LastAction().ActionChainID(),
		NPCId:         npc.ID(),
		ActionEventId: npc.ActionEvent(),
		PlayerAction:  dealAction,
	})
	if err != nil {
		return err
	}

	_, err = c.playerActions.Save(ctx, action)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinsForItemDealer) Dispose(ctx context.Context, npc *model.NPC, player *model.Player) error {
	if !player.Actions().Has[*coins_for_item.CoinsForItem](coins_for_item.PlayerActionType) {
		return nil
	}

	return c.playerActions.DeleteByNPC(ctx, npc, player)
}
