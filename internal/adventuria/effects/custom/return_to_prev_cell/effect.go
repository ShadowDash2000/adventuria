package return_to_prev_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"adventuria/pkg/mathhelper"
	"context"
	"errors"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	GetByID(ctx context.Context, id string) (*model.ActionInfo, error)
}

type boardService interface {
	MoveToCell(ctx context.Context, actor actors.Actor, cell *model.CellInfo) ([]*model.MoveResult, error)
}

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
}

var _ effects.Effect = (*ReturnToPrevCell)(nil)

const Type model.EffectType = "return_to_prev_cell"

type ReturnToPrevCell struct {
	effects.EffectBase
	actions actionsService
	board   boardService
	cells   cellsService
}

func NewDef(actions actionsService, board boardService, cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ReturnToPrevCell{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
				cells:      cells,
			}
		},
	)
}

func (r *ReturnToPrevCell) CanUse(ctx context.Context, actor *player.Actor) bool {
	if r.actions.CanDo(ctx, actor, actions.ActionTypeCompleteActivity) {
		return false
	}

	lastAction := actor.Player().LastAction()
	if lastAction == nil || lastAction.IsNew() {
		return false
	}

	rootAction, err := r.getRootAction(ctx, lastAction)
	if err != nil || rootAction.CellsPassed() == 0 {
		return false
	}

	return true
}

func (r *ReturnToPrevCell) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			lastAction := actor.Player().LastAction()

			rootAction, err := r.getRootAction(ctx, lastAction)
			if err != nil {
				return err
			}

			lastAction.SetStatus(model.ActionStatusEscape)
			lastAction, err = r.actions.Save(ctx, lastAction)
			if err != nil {
				return err
			}

			actor.Player().SetLastAction(lastAction)

			rootCell, err := r.cells.GetByID(ctx, rootAction.Cell())
			if err != nil {
				return err
			}

			cellsCount, err := r.cells.CountLocal(ctx, rootCell.World())
			if err != nil {
				return err
			}

			if cellsCount == 0 {
				return errors.New("world has no cells")
			}

			dstCellOrder := mathhelper.Mod(rootCell.LocalOrder()-rootAction.CellsPassed(), cellsCount)

			dstCell, err := r.cells.GetByLocalOrder(ctx, rootCell.World(), dstCellOrder)
			if err != nil {
				return err
			}

			_, err = r.board.MoveToCell(ctx, actor, dstCell)
			if err != nil {
				return err
			}

			actor.Player().Progress().SetCanMove(true)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

func (r *ReturnToPrevCell) getRootAction(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error) {
	rootActionId := action.RootAction()
	if rootActionId == "" {
		return action, nil
	}

	return r.actions.GetByID(ctx, rootActionId)
}
