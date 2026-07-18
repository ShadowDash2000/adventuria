package generate_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"slices"
)

type cells interface {
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.Action = (*GenerateWheel)(nil)

const Type model.ActionType = "generate_wheel"

type GenerateWheel struct {
	actions.ActionBase
	cells   cells
	actions actionsService
}

func NewDef(cells cells, actionsService actionsService) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &GenerateWheel{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				actions:    actionsService,
			}
		},
	)
}

func (g *GenerateWheel) CanDo(ctx context.Context, _ *model.Events, player *model.Player) bool {
	currentCell, err := g.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return false
	}

	if !currentCell.InCategory("activity") {
		return false
	}

	return !player.Progress().CanMove() &&
		!slices.Contains([]model.ActionType{
			actions.ActionTypeNeedToRollWheel,
			actions.ActionTypeRollWheel,
		}, player.LastAction().Type())
}

func (g *GenerateWheel) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := g.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return nil, err
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return nil, errors.New("current cell is not refreshable")
	}

	err = cellRefreshable.RefreshItems(ctx, events, player)
	if err != nil {
		return nil, err
	}

	player.LastAction().SetType(actions.ActionTypeNeedToRollWheel)

	return nil, nil
}
