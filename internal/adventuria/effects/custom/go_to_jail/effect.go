package go_to_jail

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type board interface {
	MoveToClosestCellType(ctx context.Context, events *model.Events, player *model.Player, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ model.Effect = (*GoToJail)(nil)

const Type model.EffectType = "go_to_jail"

type GoToJail struct {
	effects.EffectBase
	actions actionsService
	board   board
}

func NewDef(actions actionsService, board board) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &GoToJail{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (g *GoToJail) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	canRollWheel := g.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel)
	if canRollWheel {
		return false
	}

	canDone := g.actions.CanDo(ctx, events, player, actions.ActionTypeDone)
	canDrop := g.actions.CanDo(ctx, events, player, actions.ActionTypeDrop)

	if canDone && !canDrop {
		return false
	}

	return true
}

func (g *GoToJail) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	switch g.Value() {
	case useAfterItemAdd:
		return []event.Unsubscribe{
			events.OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
				if e.Item.Inventory().ID() != effectCtx.InvItemID {
					return e.Next()
				}

				err := g.tryToApplyEffect(ctx, events, player)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useAfterItemUse:
		return []event.Unsubscribe{
			events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
				if e.InvItemId != effectCtx.InvItemID {
					return e.Next()
				}

				err := g.tryToApplyEffect(ctx, events, player)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	}
	return nil, nil
}

func (g *GoToJail) tryToApplyEffect(ctx context.Context, events *model.Events, player *model.Player) error {
	player.Progress().SetIsInJail(true)

	_, err := g.board.MoveToClosestCellType(ctx, events, player, cells.CellTypeJail)
	return err
}
