package complete_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

var _ model.Action = (*CompleteActivity)(nil)

const Type model.ActionType = "complete_activity"

type CompleteActivity struct {
	actions.ActionBase
	cells cells
}

func NewDef(cells cells) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &CompleteActivity{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
			}
		},
	)
}

func (c *CompleteActivity) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Type() == actions.ActionTypeRollWheel
}

func (c *CompleteActivity) Do(_ context.Context, _ *model.Events, _ *model.Player, _ model.ActionRequest) (any, error) {
	return nil, errs.ErrDontDoThat
}
