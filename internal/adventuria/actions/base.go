package actions

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"slices"
)

type ActionBase struct {
	t model.ActionType
}

func NewActionBase(t model.ActionType) ActionBase {
	return ActionBase{t: t}
}

func (a ActionBase) Type() model.ActionType {
	return a.t
}

func (a ActionBase) Categories() []string {
	return Categories(a.t)
}

func (a ActionBase) InCategory(category string) bool {
	return slices.Contains(a.Categories(), category)
}

func (a ActionBase) InCategories(categories []string) bool {
	return helper.SliceContainsAll(a.Categories(), categories)
}
