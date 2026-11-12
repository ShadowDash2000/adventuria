package roll_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type cellsService interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (cells.Cell, error)
}

type activities interface {
	GetDetailedByIDs(ctx context.Context, ids []string) ([]*model.ActivityDetailed, error)
}

var _ actions.Action = (*RollWheel)(nil)

const Type model.ActionType = "roll_wheel"

type RollWheel struct {
	actions.ActionBase
	cells      cellsService
	activities activities
}

func NewDef(cells cellsService, activities activities) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &RollWheel{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				activities: activities,
			}
		},
	)
}

func (r *RollWheel) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().LastAction().Status() == model.ActionStatusNeedToRollWheel
}

func (r *RollWheel) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	currentCell, err := r.cells.GetByPlayerWrapped(ctx, actor.Player())
	if err != nil {
		return nil, err
	}

	cellRollable, ok := currentCell.(cells.Rollable)
	if !ok {
		return nil, errors.New("current cell is not rollable")
	}

	res, err := cellRollable.Roll(ctx, actor.Events(), actor.Player())
	if err != nil {
		return nil, err
	}

	lastAction := actor.Player().LastAction()
	lastAction.SetStatus(model.ActionStatusRollWheel)
	lastAction.SetActivity(res.WinnerId)

	err = actor.Events().OnAfterWheelRoll().Trigger(ctx, &model.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
