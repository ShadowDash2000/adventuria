package change_game_by_id

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type activities interface {
	GetByID(ctx context.Context, id string) (*model.Activity, error)
}

var _ model.Effect = (*ChangeGameById)(nil)

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
		func(effect model.EffectInfo) model.Effect {
			return &ChangeGameById{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (c *ChangeGameById) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	currentCell, err := c.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if ok := currentCell.Data().IsChangeGameNotAllowed(); ok {
		return false
	}

	if ok := c.actions.CanDo(ctx, events, player, actions.ActionTypeDrop); !ok {
		return false
	}

	if ok := c.actions.CanDo(ctx, events, player, actions.ActionTypeDone); !ok {
		return false
	}

	return true
}

func (c *ChangeGameById) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if e.InvItemId == effectCtx.InvItemID {
				player.LastAction().SetActivity(c.Value())
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
		events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
