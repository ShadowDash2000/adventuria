package jail_escape

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*JailEscape)(nil)

const Type model.EffectType = "jail_escape"

type JailEscape struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &JailEscape{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (j *JailEscape) CanUse(_ context.Context, actor *player.Actor) bool {
	return actor.Player().Progress().IsInJail()
}

func (j *JailEscape) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			progress := actor.Player().Progress()
			progress.SetCanMove(true)
			progress.SetIsInJail(false)
			progress.SetDropsInARow(0)

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
