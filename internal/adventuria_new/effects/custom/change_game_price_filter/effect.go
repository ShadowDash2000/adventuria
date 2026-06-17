package change_game_price_filter

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type activityFilters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilter, error)
}

var _ model.Effect = (*ChangeGamePriceFilter)(nil)

const Type model.EffectType = "change_game_price_filter"

type ChangeGamePriceFilter struct {
	effects.EffectBase
	actions         actionsService
	cells           cellsService
	activityFilters activityFilters
}

func NewDef(actions actionsService, cells cellsService, activityFilters activityFilters) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ChangeGamePriceFilter{
				EffectBase:      effects.NewEffectBase(effect),
				actions:         actions,
				cells:           cells,
				activityFilters: activityFilters,
			}
		},
	)
}

func (c *ChangeGamePriceFilter) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !c.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	currentCell, err := c.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if currentCell.Data().Type() != cells.CellTypeGame {
		return false
	}

	if currentCell.Data().IsCustomFilterNotAllowed() {
		return false
	}

	if filterId := currentCell.Data().Filter(); filterId != "" {
		filter, err := c.activityFilters.GetByID(ctx, filterId)
		if err != nil {
			return false
		}

		if len(filter.Activities()) > 0 {
			return false
		}
	}

	return true
}

func (c *ChangeGamePriceFilter) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	effectValue, err := c.decodeValue(c.Value())
	if err != nil {
		return nil, err
	}

	switch effectValue.UseType {
	case useTypeUsable:
		return []event_new.Unsubscribe{
			events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
				if e.InvItemId != effectCtx.InvItemID {
					return e.Next()
				}

				err := c.tryToApplyEffect(ctx, events, player, effectValue)
				if err != nil {
					return err
				}

				callback(ctx)

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useTypeUnusable:
		return []event_new.Unsubscribe{
			events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
				if !c.CanUse(ctx, events, player) {
					return e.Next()
				}

				err := c.tryToApplyEffect(ctx, events, player, effectValue)
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

				err := c.tryToApplyEffect(ctx, events, player, effectValue)
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

func (c *ChangeGamePriceFilter) tryToApplyEffect(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	effectValue *effectValue,
) error {
	currentCell, err := c.cells.GetCurrentCellByProgress(context.Background(), player.Progress())
	if err != nil {
		return err
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return errors.New("current cell is not refreshable")
	}

	filter := player.LastAction().CustomActivityFilter()
	if effectValue.PriceType == priceTypeMin {
		filter.MinPrice = effectValue.Price
		filter.MaxPrice = -1
	} else if effectValue.PriceType == priceTypeMax {
		filter.MinPrice = -1
		filter.MaxPrice = effectValue.Price
	}
	player.LastAction().SetCustomActivityFilter(filter)

	return cellRefreshable.RefreshItems(ctx, events, player)
}
