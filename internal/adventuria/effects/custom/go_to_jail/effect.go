package go_to_jail

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type boardService interface {
	MoveToClosestCellType(ctx context.Context, actor actors.Actor, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ effects.Effect = (*GoToJail)(nil)

const Type model.EffectType = "go_to_jail"

type GoToJail struct {
	effects.EffectBase
	actions actionsService
	board   boardService
}

func NewDef(actions actionsService, board boardService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &GoToJail{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				board:      board,
			}
		},
	)
}

func (g *GoToJail) CanUse(ctx context.Context, actor *player.Actor) bool {
	canRollWheel := g.actions.CanDo(ctx, actor, actions.ActionTypeRollWheel)
	if canRollWheel {
		return false
	}

	canDone := g.actions.CanDo(ctx, actor, actions.ActionTypeDone)
	canDrop := g.actions.CanDo(ctx, actor, actions.ActionTypeDrop)

	if canDone && !canDrop {
		return false
	}

	return true
}

func (g *GoToJail) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	switch g.Value() {
	case useAfterItemAdd:
		return []event.Unsubscribe{
			actor.Events().OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
				if e.Item.Inventory().ID() != effectCtx.InvItemID {
					return e.Next()
				}

				err := g.tryToApplyEffect(ctx, actor)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useAfterItemUse:
		return []event.Unsubscribe{
			actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
				err := g.tryToApplyEffect(ctx, actor)
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

func (g *GoToJail) tryToApplyEffect(ctx context.Context, actor *player.Actor) error {
	actor.Player().Progress().SetIsInJail(true)

	_, err := g.board.MoveToClosestCellType(ctx, actor, cells.CellTypeJail)
	return err
}
