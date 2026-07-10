package choose_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
	"slices"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type activities interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error)
}

var _ model.Effect = (*ChooseActivity)(nil)

const Type model.EffectType = "choose_activity"

type ChooseActivity struct {
	effects.EffectBase
	actions    actionsService
	activities activities
}

func NewDef(actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ChooseActivity{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (c *ChooseActivity) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	return c.actions.CanDo(ctx, events, player, actions.ActionTypeDone)
}

func (c *ChooseActivity) Subscribe(
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

			activityId, ok := e.Data["activity_id"].(string)
			if !ok {
				return errors.New("invalid activity_id")
			}

			if !slices.Contains(player.LastAction().ItemsList(), activityId) {
				return errors.New("activity_id not found in items list")
			}

			player.LastAction().SetActivity(activityId)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
