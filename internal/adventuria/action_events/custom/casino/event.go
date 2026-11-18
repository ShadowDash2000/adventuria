package casino

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

type playerActions interface {
	Save(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
	DeleteByNPCActionEvent(ctx context.Context, npc *model.NPC, player *model.Player) error
}

var _ action_events.ActionEvent = (*Casino)(nil)

const Type model.ActionEventType = "casino"

type Casino struct {
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

func NewFromJSON(items items, playerActions playerActions, data []byte) (*Casino, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &Casino{
		items:         items,
		payload:       payload,
		playerActions: playerActions,
	}, nil
}

func (c *Casino) Init(ctx context.Context, npc *model.NPC, player *model.Player) error {
	shopAction, err := shop.New(shop.Payload{
		Type:            shop.TypeCasino,
		Ids:             c.payload.ItemIds,
		PriceMultiplier: c.payload.PriceMultiplier,
		Filter: shop.Filter{
			Ids: c.payload.ItemIds,
		},
	})
	if err != nil {
		return err
	}

	action, err := model.NewPlayerAction(model.PlayerActionCreate{
		PlayerId:      player.ID(),
		ActionId:      player.LastAction().ActionChainID(),
		ActionEventId: npc.ActionEvent(),
		PlayerAction:  shopAction,
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

func (c *Casino) Dispose(ctx context.Context, npc *model.NPC, player *model.Player) error {
	if !player.Actions().Has[*shop.Shop](shop.PlayerActionType) {
		return nil
	}

	return c.playerActions.DeleteByNPCActionEvent(ctx, npc, player)
}
