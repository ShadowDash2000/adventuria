package actions

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"time"
)

type repository interface {
	Create(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	GetLastActionByPlayerId(ctx context.Context, playerId string, timeFrom time.Time) (*model.ActionInfo, error)
}

type worlds interface {
	GetDefault(ctx context.Context) (*model.World, error)
}

type cells interface {
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
}

type actionEvents interface {
	GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error)
}

type Actions struct {
	repository   repository
	worlds       worlds
	cells        cells
	actionEvents actionEvents
}

func NewActions(repository repository, worlds worlds, cells cells, actionEvents actionEvents) *Actions {
	return &Actions{
		repository:   repository,
		worlds:       worlds,
		cells:        cells,
		actionEvents: actionEvents,
	}
}

func (a *Actions) Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error) {
	return a.repository.Save(ctx, action)
}

func (a *Actions) GetLastOrDefault(ctx context.Context, playerId string, timeFrom time.Time) (*model.ActionInfo, error) {
	action, err := a.repository.GetLastActionByPlayerId(ctx, playerId, timeFrom)
	if err == nil {
		return action, nil
	} else if !errors.Is(err, errs.ErrActionNotFound) {
		return nil, err
	}

	world, err := a.worlds.GetDefault(ctx)
	if err != nil {
		return nil, err
	}

	cell, err := a.cells.GetByLocalOrder(ctx, world.ID(), 0)
	if err != nil {
		return nil, err
	}

	action, err = model.NewAction(model.ActionCreate{
		Player: playerId,
		Cell:   cell.ID(),
		Type:   "none",
	})

	return action, nil
}

type canDoContext struct {
	actionEventInfo *model.ActionEventInfo
}

func (a *Actions) getCanDoContext(ctx context.Context, player *model.Player) (*canDoContext, error) {
	actionEventInfo, err := a.actionEvents.GetByActiveCellID(ctx, player.LastAction().Cell())
	if err != nil {
		if errors.Is(err, errs.ErrActionEventNotFound) {
			return &canDoContext{}, nil
		}

		return nil, err
	}

	return &canDoContext{
		actionEventInfo: actionEventInfo,
	}, nil
}

func (a *Actions) canDoAction(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	actionDef ActionDef,
	canDoCtx *canDoContext,
) bool {
	action := actionDef.New()

	if action.CanDo(ctx, events, player) {
		return true
	}

	if canDoCtx.actionEventInfo == nil {
		return false
	}

	if canDoCtx.actionEventInfo.ActionType() != actionDef.Type() {
		return false
	}

	actionEvent, ok := action.(model.ActionEventCompatible)
	if !ok {
		return false
	}

	return actionEvent.CanDoOnEvent(ctx, events, player)
}

func (a *Actions) CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool {
	actionDef, ok := Get(t)
	if !ok {
		return false
	}

	canDoCtx, err := a.getCanDoContext(ctx, player)
	if err != nil {
		return false
	}

	return a.canDoAction(ctx, events, player, actionDef, canDoCtx)
}

func (a *Actions) Do(ctx context.Context, events *model.Events, player *model.Player, req model.ActionRequest, t model.ActionType) (any, error) {
	if actionDef, ok := Get(t); ok {
		return actionDef.New().Do(ctx, events, player, req)
	}
	return nil, errs.ErrUnknownAction
}

func (a *Actions) AvailableActions(ctx context.Context, events *model.Events, player *model.Player) []model.ActionType {
	canDoCtx, err := a.getCanDoContext(ctx, player)
	if err != nil {
		return nil
	}

	var res []model.ActionType
	for _, actionDef := range GetAll() {
		if a.canDoAction(ctx, events, player, actionDef, canDoCtx) {
			res = append(res, actionDef.Type())
		}
	}

	return res
}

func (a *Actions) HasActionsInCategory(ctx context.Context, events *model.Events, player *model.Player, category string) bool {
	canDoCtx, err := a.getCanDoContext(ctx, player)
	if err != nil {
		return false
	}

	for _, actionDef := range GetAll() {
		action := actionDef.New()

		if !action.InCategory(category) {
			continue
		}

		if a.canDoAction(ctx, events, player, actionDef, canDoCtx) {
			return true
		}
	}

	return false
}

func (a *Actions) HasActionsInCategories(ctx context.Context, events *model.Events, player *model.Player, categories []string) bool {
	canDoCtx, err := a.getCanDoContext(ctx, player)
	if err != nil {
		return false
	}

	for _, actionDef := range GetAll() {
		action := actionDef.New()

		if !action.InCategories(categories) {
			continue
		}

		if a.canDoAction(ctx, events, player, actionDef, canDoCtx) {
			return true
		}
	}

	return false
}

func (a *Actions) GetView(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) (any, error) {
	if actionDef, ok := Get(t); ok {
		action := actionDef.New()
		if actionWithView, ok := action.(model.WithView); ok {
			return actionWithView.GetView(ctx, events, player)
		}
	}
	return nil, errs.ErrUnknownAction
}
