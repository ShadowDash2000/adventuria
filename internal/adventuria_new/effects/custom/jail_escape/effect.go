package jail_escape

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*JailEscape)(nil)

const Type model.EffectType = "jail_escape"

type JailEscape struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &JailEscape{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (j *JailEscape) CanUse(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.Progress().IsInJail()
}

func (j *JailEscape) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if e.InvItemId != effectCtx.InvItemID {
				return e.Next()
			}

			progress := player.Progress()
			progress.SetIsInJail(false)
			progress.SetDropsInARow(0)
			player.LastAction().SetCanMove(true)

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
