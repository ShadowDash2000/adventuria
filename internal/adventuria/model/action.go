package model

import (
	"context"
	"errors"
)

type ActionType string
type ActionRequest any

type Action interface {
	Type() ActionType
	Categories() []string
	InCategory(string) bool
	InCategories(categories []string) bool
	CanDo(ctx context.Context, events *Events, player *Player) bool
	Do(ctx context.Context, events *Events, player *Player, actionReq ActionRequest) (any, error)
}

type ActionEventCompatible interface {
	CanDoOnEvent(ctx context.Context, events *Events, player *Player) bool
}

type ActionData struct {
	Id          string
	Player      string
	Cell        string
	Status      ActionStatus
	Activity    string
	Review      string
	CellsPassed int
	State       ActionState
	UsedItems   []string
}

type ActionInfo struct {
	data  ActionData
	isNew bool
}

type ActionCreate struct {
	Player string
	Cell   string
	Status ActionStatus
}

func NewAction(data ActionCreate) (*ActionInfo, error) {
	if data.Player == "" {
		return nil, errors.New("action: player is empty")
	}
	if data.Cell == "" {
		return nil, errors.New("action: cell is empty")
	}
	if data.Status == "" {
		return nil, errors.New("action: status is empty")
	}

	return &ActionInfo{
		data: ActionData{
			Player: data.Player,
			Cell:   data.Cell,
			Status: data.Status,
		},
		isNew: true,
	}, nil
}

func RestoreAction(data ActionData) *ActionInfo {
	return &ActionInfo{
		data:  data,
		isNew: false,
	}
}

func (a *ActionInfo) IsNew() bool {
	return a.isNew
}

func (a *ActionInfo) ID() string {
	return a.data.Id
}

func (a *ActionInfo) Player() string {
	return a.data.Player
}

func (a *ActionInfo) Cell() string {
	return a.data.Cell
}

func (a *ActionInfo) Status() ActionStatus {
	return a.data.Status
}

func (a *ActionInfo) SetStatus(s ActionStatus) {
	a.data.Status = s
}

func (a *ActionInfo) Activity() string {
	return a.data.Activity
}

func (a *ActionInfo) SetActivity(id string) {
	a.data.Activity = id
}

func (a *ActionInfo) Review() string {
	return a.data.Review
}

func (a *ActionInfo) SetReview(id string) {
	a.data.Review = id
}

func (a *ActionInfo) CellsPassed() int {
	return a.data.CellsPassed
}

func (a *ActionInfo) SetCellsPassed(count int) {
	a.data.CellsPassed = count
}

func (a *ActionInfo) State() ActionState {
	return a.data.State.Clone()
}

func (a *ActionInfo) SetState(state ActionState) {
	a.data.State = state
}

func (a *ActionInfo) UsedItems() []string {
	return a.data.UsedItems
}

func (a *ActionInfo) AddUsedItems(items ...string) {
	a.data.UsedItems = append(a.data.UsedItems, items...)
}
