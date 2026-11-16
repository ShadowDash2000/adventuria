package change_game_by_id

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

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
	actions    actionsService
	activities activities
}

func NewDef(actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ChangeGameById{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (c *ChangeGameById) CanUse(ctx context.Context, actor *player.Actor) bool {
	state := actor.Player().LastAction().State()
	if state.ActivityCompletion == nil {
		return false
	}

	if !state.ActivityCompletion.IsChangeActivityAllowed {
		return false
	}

	ok := c.actions.CanDo(ctx, actor, actions.ActionTypeDrop)
	if !ok {
		return false
	}

	ok = c.actions.CanDo(ctx, actor, actions.ActionTypeDone)
	if !ok {
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
