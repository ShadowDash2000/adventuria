package adventuria

import (
	"errors"
	"iter"
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

func (a *Actions) CanDo(user User, t ActionType) bool {
	if action, ok := a.actions[t]; ok {
		return action.CanDo(ActionContext{User: user})
	}
	return false
}

func (a *Actions) Do(user User, req ActionRequest, t ActionType) (*ActionResult, error) {
	if action, ok := a.actions[t]; ok {
		return action.Do(ActionContext{User: user}, req)
	}
	return nil, errors.New("actions: unknown action")
}

func (a *Actions) AvailableActions(user User) iter.Seq[ActionType] {
	return func(yield func(ActionType) bool) {
		for t := range a.actions {
			if !a.CanDo(user, t) {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

func (a *Actions) HasActionsInCategory(user User, category string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{User: user}) {
			continue
		}
		if action.InCategory(category) {
			return true
		}
	}
	return false
}

func (a *Actions) HasActionsInCategories(user User, categories []string) bool {
	for _, action := range a.actions {
		if !action.CanDo(ActionContext{User: user}) {
			continue
		}
		if action.InCategories(categories) {
			return true
		}
	}
	return false
}

func (a *Actions) GetVariants(user User, t ActionType) any {
	if action, ok := a.actions[t]; ok {
		return action.GetVariants(ActionContext{User: user})
	}
	return nil
}
