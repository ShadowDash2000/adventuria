package adventuria

import (
	"database/sql"
	"errors"
	"iter"

	"github.com/pocketbase/pocketbase/core"
)

type Actions struct {
	actions map[ActionType]Action
}

func NewActions() *Actions {
	a := &Actions{
		actions: make(map[ActionType]Action, len(actionsList)),
	}

	for t := range actionsList {
		action, _ := NewActionFromType(t)
		a.actions[action.Type()] = action
	}

	return a
}

func (a *Actions) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")

		user, err := GameUsers.GetByID(AppContext{App: e.App}, userId)
		if err != nil {
			return e.Next()
		}

		user.LastAction().SetProxyRecord(e.Record)

		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")

		user, err := GameUsers.GetByID(AppContext{App: e.App}, userId)
		if err != nil {
			return e.Next()
		}

		user.LastAction().SetProxyRecord(e.Record)

		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")

		user, err := GameUsers.GetByID(AppContext{App: e.App}, userId)
		if err != nil {
			return e.Next()
		}

		record, err := fetchLastUserAction(AppContext{App: e.App}, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				user.LastAction().SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
				user.LastAction().SetType(ActionTypeNone)
				user.LastAction().SetCanMove(true)

				return e.Next()
			}

			return e.Next()
		}

		user.LastAction().SetProxyRecord(record)

		return e.Next()
	})
}

func (a *Actions) CanDo(ctx AppContext, user User, t ActionType) bool {
	if action, ok := a.actions[t]; ok {
		return action.CanDo(ActionContext{
			AppContext: ctx,
			User:       user,
		})
	}
	return false
}

func (a *Actions) Do(ctx AppContext, user User, req ActionRequest, t ActionType) (*ActionResult, error) {
	if action, ok := a.actions[t]; ok {
		return action.Do(ActionContext{
			AppContext: ctx,
			User:       user,
		}, req)
	}
	return nil, errors.New("actions: unknown action")
}

func (a *Actions) AvailableActions(ctx AppContext, user User) iter.Seq[ActionType] {
	return func(yield func(ActionType) bool) {
		for t := range a.actions {
			if !a.CanDo(ctx, user, t) {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

func (a *Actions) HasActionsInCategory(ctx AppContext, user User, category string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{
			AppContext: ctx,
			User:       user,
		}) {
			continue
		}
		if action.InCategory(category) {
			return true
		}
	}
	return false
}

func (a *Actions) HasActionsInCategories(ctx AppContext, user User, categories []string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{
			AppContext: ctx,
			User:       user,
		}) {
			continue
		}
		if action.InCategories(categories) {
			return true
		}
	}
	return false
}

func (a *Actions) GetVariants(ctx AppContext, user User, t ActionType) any {
	if action, ok := a.actions[t]; ok {
		return action.GetVariants(ActionContext{
			AppContext: ctx,
			User:       user,
		})
	}
	return nil
}
