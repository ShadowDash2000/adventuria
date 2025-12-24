package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Action interface {
	Type() ActionType
	CanDo(User) bool
	Do(User, ActionRequest) (*ActionResult, error)

	setType(ActionType)
}

type ActionRecord interface {
	core.RecordProxy

	ID() string
	User() string
	SetUser(string)
	CellId() string
	setCell(string)
	Comment() string
	SetComment(string)
	Game() string
	SetGame(string)
	Type() ActionType
	SetType(ActionType)
	DiceRoll() int
	SetDiceRoll(int)
	ItemsUsed() []string
	SetItemsUsed([]string)
	ItemsList() ([]string, error)
	SetItemsList([]string)
	CanMove() bool
	SetCanMove(bool)
	CustomGameFilter() *CustomGameFilter
	ClearCustomGameFilter()
}

type CustomGameFilter struct {
	Platforms       []string
	Developers      []string
	Publishers      []string
	Genres          []string
	Tags            []string
	MinPrice        int
	MaxPrice        int
	ReleaseDateFrom types.DateTime
	ReleaseDateTo   types.DateTime
	MinCampaignTime float64
	MaxCampaignTime float64
}

type ActionType string

const (
	ActionTypeNone ActionType = "none"
	ActionTypeMove ActionType = "move"
)

type ActionRequest map[string]any

type ActionResult struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

var actionsList = map[ActionType]ActionCreator{
	ActionTypeNone: NewAction(ActionTypeNone, &NoneAction{}),
}

type ActionCreator func() Action

func RegisterActions(actions []ActionCreator) {
	for _, actionCreator := range actions {
		action := actionCreator()
		actionsList[action.Type()] = actionCreator
	}
}

func NewAction(t ActionType, a Action) ActionCreator {
	return func() Action {
		a.setType(t)
		return a
	}
}

func IsActionTypeExist(t ActionType) bool {
	_, ok := actionsList[t]
	return ok
}
