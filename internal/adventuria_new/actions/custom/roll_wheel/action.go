package roll_wheel

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type activities interface {
	GetDetailedByIDs(ctx context.Context, ids []string) ([]*model.ActivityViewDetailed, error)
}

var _ model.Action = (*RollWheel)(nil)

const Type model.ActionType = "roll_wheel"

type RollWheel struct {
	actions.ActionBase
	cells      cells
	activities activities
}

func NewActionRollWheelDef(cells cells, activities activities) actions.ActionDef {
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

func (r *RollWheel) CanDo(ctx context.Context, _ *model.Events, player *model.Player) bool {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if !currentCell.InCategory("activity") {
		return false
	}

	return !player.LastAction().CanMove() && player.LastAction().Type() != Type
}

func (r *RollWheel) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
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
	lastAction.SetType(Type)
	lastAction.SetActivity(res.WinnerId)

	err = events.OnAfterWheelRoll().Trigger(&model.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
