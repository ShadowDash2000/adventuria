package coins_for_all

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes/custom/change_balance"
	"adventuria/pkg/event"
	"context"
	"encoding/json"
)

type players interface {
	GetAll(ctx context.Context, seasonId string) ([]*model.Player, error)
}

type outboxes interface {
	Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
}

var _ model.Effect = (*CoinsForAll)(nil)

const Type model.EffectType = "coins_for_all"

type CoinsForAll struct {
	effects.EffectBase
	players  players
	outboxes outboxes
}

func NewDef(players players, outboxes outboxes) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &CoinsForAll{
				EffectBase: effects.NewEffectBase(effect),
				players:    players,
				outboxes:   outboxes,
			}
		},
	)
}

func (c *CoinsForAll) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (c *CoinsForAll) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if e.InvItemId != effectCtx.InvItemID {
				return e.Next()
			}

			effectValue, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			err = player.Progress().BalanceChange(effectValue.CoinsForPlayer)
			if err != nil {
				return err
			}

			players, err := c.players.GetAll(ctx, player.Progress().Season())
			if err != nil {
				return err
			}

			for _, p := range players {
				if p.ID() == player.ID() {
					continue
				}

				payload, err := json.Marshal(change_balance.OutboxValue{
					ProgressId: p.Progress().ID(),
					Amount:     effectValue.CoinsForOther,
				})
				outbox, err := model.NewOutbox(model.OutBoxCreate{
					Type:    change_balance.Type,
					Payload: string(payload),
				})
				if err != nil {
					return err
				}

				_, err = c.outboxes.Save(ctx, outbox)
				if err != nil {
					return err
				}
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
