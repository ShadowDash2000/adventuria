package complete_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type activityResultCalculator interface {
	Calculate(ctx context.Context, events *model.Events, cell *model.CellInfo, mode model.EffectRunMode) (*model.ActivityCompletionResult, error)
}

var _ model.Action = (*CompleteActivity)(nil)

const Type model.ActionType = "complete_activity"

type CompleteActivity struct {
	actions.ActionBase
	cells                    cells
	activityResultCalculator activityResultCalculator
}

func NewDef(cells cells, activityResultCalculator activityResultCalculator) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &CompleteActivity{
				ActionBase:               actions.NewActionBase(Type),
				cells:                    cells,
				activityResultCalculator: activityResultCalculator,
			}
		},
	)
}

func (c *CompleteActivity) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Status() == model.ActionStatusRollWheel
}

func (c *CompleteActivity) Do(_ context.Context, _ *model.Events, _ *model.Player, _ model.ActionRequest) (any, error) {
	return nil, errs.ErrDontDoThat
}
