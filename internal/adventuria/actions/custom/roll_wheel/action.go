package roll_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type activities interface {
	GetDetailedByIDs(ctx context.Context, ids []string) ([]*model.ActivityDetailed, error)
}

var _ model.Action = (*RollWheel)(nil)

const Type model.ActionType = "roll_wheel"

type RollWheel struct {
	actions.ActionBase
	cells      cells
	activities activities
}

func NewDef(cells cells, activities activities) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RollWheel{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				activities: activities,
			}
		},
	)
}

func (r *RollWheel) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Status() == model.ActionStatusNeedToRollWheel
}

func (r *RollWheel) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := r.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return nil, err
	}

	cellRollable, ok := currentCell.(model.Rollable)
	if !ok {
		return nil, errors.New("current cell is not rollable")
	}

	res, err := cellRollable.Roll(ctx, events, player)
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusRollWheel)
	lastAction.SetActivity(res.WinnerId)

	err = events.OnAfterWheelRoll().Trigger(ctx, &model.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
