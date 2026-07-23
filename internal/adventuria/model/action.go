package model

import (
	"context"
	"errors"
	"time"
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
	Id                   string
	Player               string
	Cell                 string
	Type                 ActionType
	Activity             string
	Review               string
	CellsPassed          int
	State                ActionState
	UsedItems            []string
	CustomActivityFilter CustomActivityFilter
}

type CustomActivityFilter struct {
	Platforms       []string
	Developers      []string
	Publishers      []string
	Genres          []string
	Tags            []string
	Themes          []string
	MinPrice        int
	MaxPrice        int
	ReleaseDateFrom time.Time
	ReleaseDateTo   time.Time
	MinCampaignTime float64
	MaxCampaignTime float64
}

type ActionInfo struct {
	data  ActionData
	isNew bool
}

type ActionCreate struct {
	Player string
	Cell   string
	Type   ActionType
}

func NewAction(data ActionCreate) (*ActionInfo, error) {
	if data.Player == "" {
		return nil, errors.New("action: player is empty")
	}
	if data.Cell == "" {
		return nil, errors.New("action: cell is empty")
	}
	if data.Type == "" {
		return nil, errors.New("action: type is empty")
	}

	return &ActionInfo{
		data: ActionData{
			Player: data.Player,
			Cell:   data.Cell,
			Type:   data.Type,
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

func (a *ActionInfo) Type() ActionType {
	return a.data.Type
}

func (a *ActionInfo) SetType(t ActionType) {
	a.data.Type = t
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
	return a.data.State
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

func (a *ActionInfo) CustomActivityFilter() CustomActivityFilter {
	return a.data.CustomActivityFilter
}

func (a *ActionInfo) SetCustomActivityFilter(filter CustomActivityFilter) {
	a.data.CustomActivityFilter = filter
}
