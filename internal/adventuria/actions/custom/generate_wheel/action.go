package generate_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
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
	if player.LastAction().Type() != actions.ActionTypeRollDice {
		return false
	}

	currentCell, err := g.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	return currentCell.InCategory("activity")
}

func (g *GenerateWheel) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := g.cells.GetCurrentCellByProgress(ctx, player.Progress())
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
