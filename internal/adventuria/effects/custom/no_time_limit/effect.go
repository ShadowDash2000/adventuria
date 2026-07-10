package no_time_limit

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

var _ model.Effect = (*NoTimeLimit)(nil)

const Type model.EffectType = "no_time_limit"

type NoTimeLimit struct {
	effects.EffectBase
	actions actionsService
	cells   cellsService
}

func NewDef(actions actionsService, cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &NoTimeLimit{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				cells:      cells,
			}
		},
	)
}

func (n *NoTimeLimit) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !n.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	currentCell, err := n.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if currentCell.Data().Type() != cells.CellTypeGame {
		return false
	}

	return true
}

func (n *NoTimeLimit) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			if !n.CanUse(ctx, events, player) {
				return e.Next()
			}

			err := n.tryToApplyEffect(ctx, events, player)
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

			if !n.CanUse(ctx, events, player) {
				return e.Next()
			}

			err := n.tryToApplyEffect(ctx, events, player)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

func (n *NoTimeLimit) tryToApplyEffect(ctx context.Context, events *model.Events, player *model.Player) error {
	currentCell, err := n.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return err
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return errors.New("current cell is not refreshable")
	}

	filter := player.LastAction().CustomActivityFilter()
	filter.MinCampaignTime = -1
	filter.MaxCampaignTime = -1
	player.LastAction().SetCustomActivityFilter(filter)

	return cellRefreshable.RefreshItems(ctx, events, player)
}
