package return_to_prev_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	GetByID(ctx context.Context, id string) (*model.ActionInfo, error)
}

type board interface {
	Move(ctx context.Context, actor model.BoardActor, steps int, moveType model.MoveType) ([]*model.MoveResult, error)
}

var _ effects.Effect = (*ReturnToPrevCell)(nil)

const Type model.EffectType = "return_to_prev_cell"

type ReturnToPrevCell struct {
	effects.EffectBase
	actions actionsService
	board   board
}

func NewDef(actions actionsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ReturnToPrevCell{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (r *ReturnToPrevCell) CanUse(ctx context.Context, actor *player.Actor) bool {
	if actor.Player().Progress().CellsPassed() == 0 {
		return false
	}

	return !r.actions.CanDo(ctx, actor, actions.ActionTypeCompleteActivity)
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
			lastAction.SetStatus(model.ActionStatusEscape)
			lastAction, err := r.actions.Save(ctx, lastAction)
			if err != nil {
				return err
			}

			rootAction := lastAction
			rootActionId := rootAction.RootAction()
			if rootActionId != "" {
				action, err := r.actions.GetByID(ctx, rootActionId)
				if err != nil {
					return err
				}
				rootAction = action
			}

			_, err = r.board.Move(ctx, actor, -rootAction.CellsPassed(), model.MoveTypePath)
			if err != nil {
				return err
			}

			actor.Player().Progress().SetCanMove(true)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
