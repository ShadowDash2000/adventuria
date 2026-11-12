package coins_for_all

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes/custom/change_balance"
	"adventuria/pkg/event"
	"context"
)

type progress interface {
	GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.PlayerProgress, error)
}

type notify interface {
	BalanceChangeNotify(ctx context.Context, id string, amount int) error
}

type outboxes interface {
	Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
}

var _ effects.Effect = (*CoinsForAll)(nil)

const Type model.EffectType = "coins_for_all"

type CoinsForAll struct {
	effects.EffectBase
	progress progress
	notify   notify
	outboxes outboxes
}

func NewDef(progress progress, notify notify, outboxes outboxes) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &CoinsForAll{
				EffectBase: effects.NewEffectBase(effect),
				progress:   progress,
				notify:     notify,
				outboxes:   outboxes,
			}
		},
	)
}

func (c *CoinsForAll) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (c *CoinsForAll) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			effectValue, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			err = actor.Player().Progress().BalanceChange(effectValue.CoinsForPlayer)
			if err != nil {
				return err
			}

			err = c.notify.BalanceChangeNotify(ctx, actor.Player().Progress().ID(), effectValue.CoinsForPlayer)
			if err != nil {
				return err
			}

			progresses, err := c.progress.GetAllBySeasonID(ctx, actor.Player().Progress().Season())
			if err != nil {
				return err
			}

			for _, p := range progresses {
				if p.Player() == actor.Player().ID() {
					continue
				}

				task, err := change_balance.New(change_balance.Payload{
					ProgressId: p.ID(),
					Amount:     &effectValue.CoinsForOther,
				})
				if err != nil {
					return err
				}

				outbox, err := model.NewOutbox(model.OutboxCreate{
					Task: task,
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
