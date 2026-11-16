package complete_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type doneResultCalculator interface {
	Calculate(ctx context.Context, actor *player.Actor, mode model.EffectRunMode) (*model.ActivityDoneResult, error)
}

type dropResultCalculator interface {
	Calculate(ctx context.Context, actor *player.Actor, mode model.EffectRunMode) (*model.ActivityDropResult, error)
}

var _ actions.Action = (*CompleteActivity)(nil)

const Type model.ActionType = "complete_activity"

type CompleteActivity struct {
	actions.ActionBase
	doneResultCalculator doneResultCalculator
	dropResultCalculator dropResultCalculator
}

func NewDef(
	doneResultCalculator doneResultCalculator,
	dropResultCalculator dropResultCalculator,
) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &CompleteActivity{
				ActionBase:           actions.NewActionBase(Type),
				doneResultCalculator: doneResultCalculator,
				dropResultCalculator: dropResultCalculator,
			}
		},
	)
}

func (c *CompleteActivity) CanDo(_ context.Context, actor *player.Actor) bool {
	if actor.Player().LastAction().Activity() == "" {
		return false
	}

	return actor.Player().LastAction().Status() == model.ActionStatusRollWheel
}

func (c *CompleteActivity) Do(_ context.Context, _ *player.Actor, _ actions.Request) (any, error) {
	return nil, errs.ErrDontDoThat
}
