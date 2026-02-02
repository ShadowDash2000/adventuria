package adventuria

import (
	"errors"
	"slices"
)

type ActionBase struct {
	t    ActionType
	user User
}

func NewActionFromType(actionType ActionType) (Action, error) {
	actionDef, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionDef.New()

	return action, nil
}

func (a *ActionBase) Type() ActionType {
	return a.t
}

func (a *ActionBase) Categories() []string {
	if def, ok := actionsList[a.t]; ok {
		return def.Categories
	}

	return nil
}

func (a *ActionBase) setType(t ActionType) {
	a.t = t
}

func (a *ActionBase) InCategory(category string) bool {
	return slices.Contains(a.Categories(), category)
}

func (a *ActionBase) InCategories(categories []string) bool {
	return SliceContainsAll(a.Categories(), categories)
}
