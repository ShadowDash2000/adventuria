package cell_points_divide

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*CellPointsDivide)(nil)

const Type model.EffectType = "cell_points_divide"

type CellPointsDivide struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &CellPointsDivide{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (c *CellPointsDivide) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (c *CellPointsDivide) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnActivityDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDoneEvent) error {
			divider, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			e.Result.SetPoints(e.Result.Points() / divider)

			if e.Mode == model.EffectRunModeApply {
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
		actor.Events().OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
