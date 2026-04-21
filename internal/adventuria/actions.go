package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/result"
	"database/sql"
	"errors"
	"iter"

	"github.com/pocketbase/pocketbase/core"
)

type Actions struct {
	actions map[ActionType]Action
}

func NewActions(ctx AppContext) *Actions {
	a := &Actions{
		actions: make(map[ActionType]Action, len(actionsList)),
	}

	for t := range actionsList {
		action, _ := NewActionFromType(t)
		a.actions[action.Type()] = action
	}

	a.bindHooks(ctx)

	return a
}

func (a *Actions) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		playerId := e.Record.GetString(schema.ActionSchema.Player)

		player, err := GamePlayers.GetByID(AppContext{App: e.App}, playerId)
		if err != nil {
			return e.Next()
		}

		player.LastAction().SetProxyRecord(e.Record)

		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		playerId := e.Record.GetString(schema.ActionSchema.Player)

		player, err := GamePlayers.GetByID(AppContext{App: e.App}, playerId)
		if err != nil {
			return e.Next()
		}

		player.LastAction().SetProxyRecord(e.Record)

		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		playerId := e.Record.GetString(schema.ActionSchema.Player)

		player, err := GamePlayers.GetByID(AppContext{App: e.App}, playerId)
		if err != nil {
			return e.Next()
		}

		record, err := fetchLastPlayerAction(AppContext{App: e.App}, playerId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				player.LastAction().SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionActions)))
				player.LastAction().SetType(ActionTypeNone)
				player.LastAction().SetCanMove(true)

				return e.Next()
			}

			return e.Next()
		}

		player.LastAction().SetProxyRecord(record)

		return e.Next()
	})
}

func (a *Actions) CanDo(ctx AppContext, player Player, t ActionType) bool {
	if action, ok := a.actions[t]; ok {
		return action.CanDo(ActionContext{
			AppContext: ctx,
			Player:     player,
		})
	}
	return false
}

func (a *Actions) Do(ctx AppContext, player Player, req ActionRequest, t ActionType) (*result.Result, error) {
	if action, ok := a.actions[t]; ok {
		return action.Do(ActionContext{
			AppContext: ctx,
			Player:     player,
		}, req)
	}
	return result.Err("unknown action"), errors.New("actions: unknown action")
}

func (a *Actions) AvailableActions(ctx AppContext, player Player) iter.Seq[ActionType] {
	return func(yield func(ActionType) bool) {
		for t := range a.actions {
			if !a.CanDo(ctx, player, t) {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

func (a *Actions) HasActionsInCategory(ctx AppContext, player Player, category string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{
			AppContext: ctx,
			Player:     player,
		}) {
			continue
		}
		if action.InCategory(category) {
			return true
		}
	}
	return false
}

func (a *Actions) HasActionsInCategories(ctx AppContext, player Player, categories []string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{
			AppContext: ctx,
			Player:     player,
		}) {
			continue
		}
		if action.InCategories(categories) {
			return true
		}
	}
	return false
}

func (a *Actions) GetVariants(ctx AppContext, player Player, t ActionType) any {
	if action, ok := a.actions[t]; ok {
		return action.GetVariants(ActionContext{
			AppContext: ctx,
			Player:     player,
		})
	}
	return nil
}
