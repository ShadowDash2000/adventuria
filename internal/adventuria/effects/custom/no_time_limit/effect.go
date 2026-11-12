package no_time_limit

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ effects.Effect = (*NoTimeLimit)(nil)

const Type model.EffectType = "no_time_limit"

type NoTimeLimit struct {
	effects.EffectBase
	actions    actionsService
	activities activities
}

func NewDef(actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &NoTimeLimit{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (n *NoTimeLimit) CanUse(ctx context.Context, actor *player.Actor) bool {
	if !n.actions.CanDo(ctx, actor, actions.ActionTypeRollWheel) {
		return false
	}

	activityFilter := actor.Player().LastAction().State().ActivityFilter
	if activityFilter == nil {
		return false
	}

	if activityFilter.Type != model.ActivityTypeGame {
		return false
	}

	return true
}

func (n *NoTimeLimit) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			if !n.CanUse(ctx, actor) {
				return e.Next()
			}

			err := n.tryToApplyEffect(ctx, actor.Player())
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
		actor.Events().OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
			if e.Item.Inventory().ID() != effectCtx.InvItemID {
				return e.Next()
			}

			if !n.CanUse(ctx, actor) {
				return e.Next()
			}

			err := n.tryToApplyEffect(ctx, actor.Player())
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

func (n *NoTimeLimit) tryToApplyEffect(ctx context.Context, player *model.Player) error {
	actionState := player.LastAction().State()
	if actionState.ActivityFilter == nil {
		return errs.ErrNoActiveActivityFilter
	}

	actionState.ActivityFilter.MinCampaignTime = -1
	actionState.ActivityFilter.MaxCampaignTime = -1

	ids, err := n.activities.GetRandomIDsByFilter(ctx, actionState.ActivityFilter.Clone())
	if err != nil {
		return err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}

	return nil
}
