package choose_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
	"slices"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type activities interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error)
}

var _ effects.Effect = (*ChooseActivity)(nil)

const Type model.EffectType = "choose_activity"

type ChooseActivity struct {
	effects.EffectBase
	actions    actionsService
	activities activities
}

func NewDef(actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ChooseActivity{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (c *ChooseActivity) CanUse(ctx context.Context, actor *player.Actor) bool {
	activitiesState := actor.Player().LastAction().State().Activities

	if activitiesState == nil {
		return false
	}

	if len(activitiesState.Ids) == 0 {
		return false
	}

	return c.actions.CanDo(ctx, actor, actions.ActionTypeCompleteActivity)
}

func (c *ChooseActivity) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			activityId, ok := e.Data["activity_id"].(string)
			if !ok {
				return errors.New("invalid activity_id")
			}

			activitiesState := actor.Player().LastAction().State().Activities
			if activitiesState == nil {
				return errs.ErrNoActiveActivity
			}

			if !slices.Contains(activitiesState.Ids, activityId) {
				return errors.New("activity_id not found in items list")
			}

			actor.Player().LastAction().SetActivity(activityId)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
