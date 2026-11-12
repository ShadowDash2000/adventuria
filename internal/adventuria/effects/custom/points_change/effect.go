package points_change

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*PointsChange)(nil)

const Type model.EffectType = "points_change"

type PointsChange struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &PointsChange{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (p *PointsChange) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (p *PointsChange) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnActivityDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDoneEvent) error {
			amount, err := p.decodeValue(p.Value())
			if err != nil {
				return err
			}

			e.Result.SetPoints(e.Result.Points() + amount)

			if e.Mode == model.EffectRunModeApply {
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
