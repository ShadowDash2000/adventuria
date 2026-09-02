package change_game_price_filter

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cellsService interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ model.Effect = (*ChangeGamePriceFilter)(nil)

const Type model.EffectType = "change_game_price_filter"

type ChangeGamePriceFilter struct {
	effects.EffectBase
	actions    actionsService
	cells      cellsService
	activities activities
}

func NewDef(actions actionsService, cells cellsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ChangeGamePriceFilter{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				cells:      cells,
				activities: activities,
			}
		},
	)
}

func (c *ChangeGamePriceFilter) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !c.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	currentCell, err := c.cells.GetByPlayer(ctx, player)
	if err != nil {
		return false
	}

	if currentCell.IsCustomFilterNotAllowed() {
		return false
	}

	activityFilter := player.LastAction().State().ActivityFilter
	if activityFilter == nil {
		return false
	}

	if activityFilter.Type != model.ActivityTypeGame {
		return false
	}

	if len(activityFilter.Activities) > 0 {
		return false
	}

	return true
}

func (c *ChangeGamePriceFilter) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	effectValue, err := c.decodeValue(c.Value())
	if err != nil {
		return nil, err
	}

	switch effectValue.UseType {
	case useTypeUsable:
		return []event.Unsubscribe{
			events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
				if e.InvItemId != effectCtx.InvItemID {
					return e.Next()
				}

				err := c.tryToApplyEffect(ctx, player, effectValue)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useTypeUnusable:
		return []event.Unsubscribe{
			events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
				if !c.CanUse(ctx, events, player) {
					return e.Next()
				}

				err := c.tryToApplyEffect(ctx, player, effectValue)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
			events.OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
				if e.Item.Inventory().ID() != effectCtx.InvItemID {
					return e.Next()
				}

				if !c.CanUse(ctx, events, player) {
					return e.Next()
				}

				err := c.tryToApplyEffect(ctx, player, effectValue)
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

func (c *ChangeGamePriceFilter) tryToApplyEffect(ctx context.Context, player *model.Player, effectValue *effectValue) error {
	actionState := player.LastAction().State()
	if actionState.ActivityFilter == nil {
		return errs.ErrNoActiveActivityFilter
	}

	if effectValue.PriceType == priceTypeMin {
		actionState.ActivityFilter.MinPrice = effectValue.Price
		actionState.ActivityFilter.MaxPrice = -1
	} else if effectValue.PriceType == priceTypeMax {
		actionState.ActivityFilter.MinPrice = -1
		actionState.ActivityFilter.MaxPrice = effectValue.Price
	}

	ids, err := c.activities.GetRandomIDsByFilter(ctx, *actionState.ActivityFilter)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return errs.ErrInvalidActivityFilter
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}

	return nil
}
