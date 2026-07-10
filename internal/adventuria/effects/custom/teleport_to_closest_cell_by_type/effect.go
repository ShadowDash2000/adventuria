package teleport_to_closest_cell_by_type

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
	HasActionsInCategories(ctx context.Context, events *model.Events, player *model.Player, categories []string) bool
}

type board interface {
	MoveToClosestCellType(ctx context.Context, events *model.Events, player *model.Player, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ model.Effect = (*TeleportToClosestCellByType)(nil)

const Type model.EffectType = "teleport_to_closest_cell_by_type"

type TeleportToClosestCellByType struct {
	effects.EffectBase
	actions actionsService
	board   board
}

func NewDef(actions actionsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &TeleportToClosestCellByType{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (t *TeleportToClosestCellByType) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if t.actions.HasActionsInCategories(ctx, events, player, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	canDone := t.actions.CanDo(ctx, events, player, actions.ActionTypeDone)
	canDrop := t.actions.CanDo(ctx, events, player, actions.ActionTypeDrop)

	if canDone && !canDrop {
		return false
	}

	return true
}

func (t *TeleportToClosestCellByType) Subscribe(
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

			cellType := t.decodeValue(t.Value())
			_, err := t.board.MoveToClosestCellType(ctx, events, player, cellType)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
