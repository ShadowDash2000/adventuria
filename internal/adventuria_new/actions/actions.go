package actions

import (
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

type Actions struct {
	repository repository
	worlds     worlds
	cells      cells
}

func NewActions(repository repository, worlds worlds, cells cells) *Actions {
	return &Actions{
		repository: repository,
		worlds:     worlds,
		cells:      cells,
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

	action, err = model.NewAction(uuid.New(), model.ActionCreate{
		Player: playerId,
		Cell:   cell.ID(),
		Type:   "none",
	})

	return action, nil
}

func (a *Actions) CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool {
	if actionDef, ok := Get(t); ok {
		return actionDef.New().CanDo(ctx, events, player)
	}
	return false
}

func (a *Actions) Do(ctx context.Context, events *model.Events, player *model.Player, req model.ActionRequest, t model.ActionType) (any, error) {
	if actionDef, ok := Get(t); ok {
		return actionDef.New().Do(ctx, events, player, req)
	}
	return nil, errs.ErrUnknownAction
}

func (a *Actions) AvailableActions(ctx context.Context, events *model.Events, player *model.Player) []model.ActionType {
	var res []model.ActionType
	for _, actionDef := range GetAll() {
		if actionDef.New().CanDo(ctx, events, player) {
			res = append(res, actionDef.Type())
		}
	}
	return res
}

func (a *Actions) HasActionsInCategory(ctx context.Context, events *model.Events, player *model.Player, category string) bool {
	for _, actionDef := range GetAll() {
		action := actionDef.New()
		if !action.CanDo(ctx, events, player) {
			continue
		}
		if action.InCategory(category) {
			return true
		}
	}
	return false
}

func (a *Actions) HasActionsInCategories(ctx context.Context, events *model.Events, player *model.Player, categories []string) bool {
	for _, actionDef := range GetAll() {
		action := actionDef.New()
		if !action.CanDo(ctx, events, player) {
			continue
		}
		if action.InCategories(categories) {
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
