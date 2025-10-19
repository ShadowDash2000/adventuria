package adventuria

import (
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type Action interface {
	core.RecordProxy

	ID() string
	Save() error
	User() User
	UserId() string
	CellId() string
	SetCell(string)
	Comment() string
	SetComment(string)
	Game() string
	SetGame(string)
	Type() ActionType
	SetType(ActionType)
	SetNotAffectNextStep(bool)
	DiceRoll() int
	SetDiceRoll(int)
	ItemsUsed() []string
	SetItemsUsed([]string)
	Seed() int
	SetSeed(int)

	CanDo() bool
	NextAction() ActionType
	Do(ActionRequest) (*ActionResult, error)

	setUser(user User)
}

type ActionType string

type ActionRequest struct {
	Comment string
}

type ActionResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var actionsList = map[ActionType]ActionCreator{
	"none": NewAction(&NoneAction{}),
}

type ActionCreator func() Action

func RegisterActions(actions map[ActionType]ActionCreator) {
	maps.Insert(actionsList, maps.All(actions))
}

func NewAction(a Action) ActionCreator {
	return func() Action {
		return a
	}
}
