package return_to_prev_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type board interface {
	Move(ctx context.Context, events *model.Events, player *model.Player, steps int) ([]*model.MoveResult, error)
}

var _ model.Effect = (*ReturnToPrevCell)(nil)

const Type model.EffectType = "return_to_prev_cell"

type ReturnToPrevCell struct {
	effects.EffectBase
	actions actionsService
	board   board
}

func NewDef(actions actionsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ReturnToPrevCell{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (r *ReturnToPrevCell) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if player.Progress().CellsPassed() == 0 {
		return false
	}

	return !r.actions.CanDo(ctx, events, player, actions.ActionTypeDone)
}

func (r *ReturnToPrevCell) Subscribe(
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

			_, err := r.board.Move(ctx, events, player, -player.LastAction().CellsPassed())
			if err != nil {
				return err
			}

			player.Progress().SetCanMove(true)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
