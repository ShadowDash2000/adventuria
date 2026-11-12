package complete_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type doneResultCalculator interface {
	Calculate(ctx context.Context, events *model.Events, cell *model.CellInfo, mode model.EffectRunMode) (*model.ActivityDoneResult, error)
}

type dropResultCalculator interface {
	Calculate(ctx context.Context, events *model.Events, cell *model.CellInfo, mode model.EffectRunMode) (*model.ActivityDropResult, error)
}

var _ actions.Action = (*CompleteActivity)(nil)

const Type model.ActionType = "complete_activity"

type CompleteActivity struct {
	actions.ActionBase
	cells                cells
	doneResultCalculator doneResultCalculator
	dropResultCalculator dropResultCalculator
}

func NewDef(
	cells cells,
	doneResultCalculator doneResultCalculator,
	dropResultCalculator dropResultCalculator,
) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &CompleteActivity{
				ActionBase:           actions.NewActionBase(Type),
				cells:                cells,
				doneResultCalculator: doneResultCalculator,
				dropResultCalculator: dropResultCalculator,
			}
		},
	)
}

func (c *CompleteActivity) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().LastAction().Status() == model.ActionStatusRollWheel
}

func (c *CompleteActivity) Do(_ context.Context, _ *player.Actor, _ actions.Request) (any, error) {
	return nil, errs.ErrDontDoThat
}
