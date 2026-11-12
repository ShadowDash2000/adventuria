package change_game_by_id

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type activities interface {
	GetByID(ctx context.Context, id string) (*model.Activity, error)
}

var _ effects.Effect = (*ChangeGameById)(nil)

const Type model.EffectType = "change_game_by_id"

type ChangeGameById struct {
	effects.EffectBase
	cells      cells
	actions    actionsService
	activities activities
}

func NewDef(cells cells, actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ChangeGameById{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (c *ChangeGameById) CanUse(ctx context.Context, actor *player.Actor) bool {
	currentCell, err := c.cells.GetByPlayer(ctx, actor.Player())
	if err != nil {
		return false
	}

	if ok := currentCell.IsChangeGameNotAllowed(); ok {
		return false
	}

	if ok := c.actions.CanDo(ctx, actor, actions.ActionTypeDrop); !ok {
		return false
	}

	if ok := c.actions.CanDo(ctx, actor, actions.ActionTypeDone); !ok {
		return false
	}

	return true
}

func (c *ChangeGameById) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			actor.Player().LastAction().SetActivity(c.Value())
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
		actor.Events().OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
