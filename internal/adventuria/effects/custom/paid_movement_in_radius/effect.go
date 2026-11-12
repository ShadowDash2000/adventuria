package paid_movement_in_radius

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	HasActionsInCategories(ctx context.Context, actor *player.Actor, categories []string) bool
}

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
}

type boardService interface {
	GetLocalDistanceBetweenCells(cell1, cell2 *model.CellInfo) (int, error)
	MoveToCellId(ctx context.Context, actor actors.Actor, cellId string) ([]*model.MoveResult, error)
}

var _ effects.Effect = (*PaidMovementInRadius)(nil)

const Type model.EffectType = "paid_movement_in_radius"

type PaidMovementInRadius struct {
	effects.EffectBase
	actions actionsService
	cells   cellsService
	board   boardService
}

func NewDef(actions actionsService, cells cellsService, board boardService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &PaidMovementInRadius{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				cells:      cells,
				board:      board,
			}
		},
	)
}

func (p *PaidMovementInRadius) CanUse(ctx context.Context, actor *player.Actor) bool {
	effectValue, err := p.decodeValue(p.Value())
	if err != nil {
		return false
	}

	if actor.Player().Progress().Balance() < effectValue.Price {
		return false
	}

	if p.actions.CanDo(ctx, actor, actions.ActionTypeRollDice) {
		return true
	}

	if p.actions.HasActionsInCategories(ctx, actor, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, err := p.cells.GetByID(ctx, actor.Player().LastAction().Cell())
	if err != nil {
		return false
	}

	canDone := p.actions.CanDo(ctx, actor, actions.ActionTypeDone)
	canDrop := p.actions.CanDo(ctx, actor, actions.ActionTypeDrop)

	if canDone && !canDrop {
		if currentCell.Type() != cells.CellTypeJail {
			return false
		}
	}

	return true
}

func (p *PaidMovementInRadius) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			cellId, ok := e.Data["cell_id"].(string)
			if !ok {
				return errors.New("invalid cell_id")
			}

			effectValue, err := p.decodeValue(p.Value())
			if err != nil {
				return err
			}

			currentCell, err := p.cells.GetByID(ctx, actor.Player().LastAction().Cell())
			if err != nil {
				return err
			}

			destinationCell, err := p.cells.GetByID(ctx, cellId)
			if err != nil {
				return err
			}

			distance, err := p.board.GetLocalDistanceBetweenCells(currentCell, destinationCell)
			if err != nil {
				return err
			}

			if distance == 0 || distance > effectValue.Radius {
				return errors.New("destination cell is too far")
			}

			_, err = p.board.MoveToCellId(ctx, actor, cellId)
			if err != nil {
				return err
			}

			err = actor.Player().Progress().BalanceChange(-effectValue.Price)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
