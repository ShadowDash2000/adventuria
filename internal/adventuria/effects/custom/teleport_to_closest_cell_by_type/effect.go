package teleport_to_closest_cell_by_type

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	HasActionsInCategories(ctx context.Context, actor *player.Actor, categories []string) bool
}

type boardService interface {
	MoveToClosestCellType(ctx context.Context, actor actors.Actor, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ effects.Effect = (*TeleportToClosestCellByType)(nil)

const Type model.EffectType = "teleport_to_closest_cell_by_type"

type TeleportToClosestCellByType struct {
	effects.EffectBase
	actions actionsService
	board   boardService
}

func NewDef(actions actionsService, board boardService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &TeleportToClosestCellByType{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (t *TeleportToClosestCellByType) CanUse(ctx context.Context, actor *player.Actor) bool {
	if t.actions.HasActionsInCategories(ctx, actor, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	canDone := t.actions.CanDo(ctx, actor, actions.ActionTypeDone)
	canDrop := t.actions.CanDo(ctx, actor, actions.ActionTypeDrop)

	if canDone && !canDrop {
		return false
	}

	return true
}

func (t *TeleportToClosestCellByType) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			cellType := t.decodeValue(t.Value())
			_, err := t.board.MoveToClosestCellType(ctx, actor, cellType)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
