package paid_movement_in_radius

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
	HasActionsInCategories(ctx context.Context, events *model.Events, player *model.Player, categories []string) bool
}

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
}

type board interface {
	GetLocalDistanceBetweenCells(cell1, cell2 *model.CellInfo) (int, error)
	MoveToCellId(ctx context.Context, events *model.Events, player *model.Player, cellId string) ([]*model.MoveResult, error)
}

var _ model.Effect = (*PaidMovementInRadius)(nil)

const Type model.EffectType = "paid_movement_in_radius"

type PaidMovementInRadius struct {
	effects.EffectBase
	actions actionsService
	cells   cellsService
	board   board
}

func NewDef(actions actionsService, cells cellsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &PaidMovementInRadius{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				cells:      cells,
				board:      board,
			}
		},
	)
}

func (p *PaidMovementInRadius) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	effectValue, err := p.decodeValue(p.Value())
	if err != nil {
		return false
	}

	if player.Progress().Balance() < effectValue.Price {
		return false
	}

	if p.actions.CanDo(ctx, events, player, actions.ActionTypeRollDice) {
		return true
	}

	if p.actions.HasActionsInCategories(ctx, events, player, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, err := p.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	canDone := p.actions.CanDo(ctx, events, player, actions.ActionTypeDone)
	canDrop := p.actions.CanDo(ctx, events, player, actions.ActionTypeDrop)

	if canDone && !canDrop {
		if currentCell.Data().Type() != cells.CellTypeJail {
			return false
		}
	}

	return true
}

func (p *PaidMovementInRadius) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if e.InvItemId != effectCtx.InvItemID {
				return e.Next()
			}

			cellId, ok := e.Data["cell_id"].(string)
			if !ok {
				return errors.New("invalid cell_id")
			}

			effectValue, err := p.decodeValue(p.Value())
			if err != nil {
				return err
			}

			currentCell, err := p.cells.GetCurrentCellByProgress(ctx, player.Progress())
			if err != nil {
				return err
			}

			destinationCell, err := p.cells.GetByID(ctx, cellId)
			if err != nil {
				return err
			}

			distance, err := p.board.GetLocalDistanceBetweenCells(currentCell.Data(), destinationCell)
			if err != nil {
				return err
			}

			if distance == 0 || distance > effectValue.Radius {
				return errors.New("destination cell is too far")
			}

			_, err = p.board.MoveToCellId(ctx, events, player, cellId)
			if err != nil {
				return err
			}

			err = player.Progress().BalanceChange(-effectValue.Price)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
