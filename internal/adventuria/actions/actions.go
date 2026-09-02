package actions

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	Update(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	GetLastPlayerActionBySeasonID(ctx context.Context, playerId, seasonId string) (*model.ActionInfo, error)
	GetByID(ctx context.Context, id string) (*model.ActionInfo, error)
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
	if action.IsNew() {
		return a.repository.Create(ctx, action)
	}

	return a.repository.Update(ctx, action)
}

func (a *Actions) CreateAndSetLastAction(ctx context.Context, player *model.Player, action *model.ActionInfo) (*model.ActionInfo, error) {
	if !action.IsNew() {
		return nil, errors.New("action must be new")
	}

	action, err := a.Save(ctx, action)
	if err != nil {
		return nil, err
	}

	player.SetLastAction(action)
	player.Progress().SetLastAction(action.ID())

	return action, nil
}

func (a *Actions) GetLastOrDefault(ctx context.Context, playerId, seasonId string) (*model.ActionInfo, error) {
	action, err := a.repository.GetLastPlayerActionBySeasonID(ctx, playerId, seasonId)
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
		Status: model.ActionStatusNone,
	})

	return action, nil
}

func (a *Actions) canDoAction(ctx context.Context, events *model.Events, player *model.Player, actionDef ActionDef) bool {
	action := actionDef.New()

	return action.CanDo(ctx, events, player)
}

func (a *Actions) CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool {
	actionDef, ok := Get(t)
	if !ok {
		return false
	}

	return a.canDoAction(ctx, events, player, actionDef)
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
		if a.canDoAction(ctx, events, player, actionDef) {
			res = append(res, actionDef.Type())
		}
	}

	return res
}

func (a *Actions) HasActionsInCategory(ctx context.Context, events *model.Events, player *model.Player, category string) bool {
	for _, actionDef := range GetAll() {
		action := actionDef.New()

		if !action.InCategory(category) {
			continue
		}

		if a.canDoAction(ctx, events, player, actionDef) {
			return true
		}
	}

	return false
}

func (a *Actions) HasActionsInCategories(ctx context.Context, events *model.Events, player *model.Player, categories []string) bool {
	for _, actionDef := range GetAll() {
		action := actionDef.New()

		if !action.InCategories(categories) {
			continue
		}

		if a.canDoAction(ctx, events, player, actionDef) {
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

func (a *Actions) GetByID(ctx context.Context, id string) (*model.ActionInfo, error) {
	return a.repository.GetByID(ctx, id)
}
