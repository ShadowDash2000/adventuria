package teleport_to_random_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
	"math/rand"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	HasActionsInCategories(ctx context.Context, actor *player.Actor, categories []string) bool
}

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error)
	CountGlobal(ctx context.Context) (int, error)
}

type board interface {
	MoveToCellId(ctx context.Context, actor model.BoardActor, cellId string) ([]*model.MoveResult, error)
}

var _ effects.Effect = (*TeleportToRandomCell)(nil)

const Type model.EffectType = "teleport_to_random_cell"

type TeleportToRandomCell struct {
	effects.EffectBase
	actions actionsService
	cells   cellsService
	board   board
}

func NewDef(actions actionsService, cells cellsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &TeleportToRandomCell{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				cells:      cells,
				board:      board,
			}
		},
	)
}

func (t *TeleportToRandomCell) CanUse(ctx context.Context, actor *player.Actor) bool {
	if t.actions.CanDo(ctx, actor, actions.ActionTypeRollDice) {
		return true
	}

	if t.actions.HasActionsInCategories(ctx, actor, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, err := t.cells.GetByID(ctx, actor.Player().LastAction().Cell())
	if err != nil {
		return false
	}

	canDone := t.actions.CanDo(ctx, actor, actions.ActionTypeDone)
	canDrop := t.actions.CanDo(ctx, actor, actions.ActionTypeDrop)

	if canDone && !canDrop {
		if currentCell.Type() != cells.CellTypeJail {
			return false
		}
	}

	return true
}

func (t *TeleportToRandomCell) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			currentCell, err := t.cells.GetByID(ctx, actor.Player().LastAction().Cell())
			if err != nil {
				return err
			}

			cellsCountGlobal, err := t.cells.CountGlobal(ctx)
			if err != nil {
				return err
			}

			if cellsCountGlobal <= 1 {
				return errors.New("not enough cells to teleport to")
			}

			randomCellOrder := rand.Intn(cellsCountGlobal - 1)
			if randomCellOrder >= currentCell.GlobalOrder() {
				randomCellOrder++
			}

			randomCell, err := t.cells.GetByGlobalOrder(ctx, randomCellOrder)
			if err != nil {
				return err
			}

			_, err = t.board.MoveToCellId(ctx, actor, randomCell.ID())
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
