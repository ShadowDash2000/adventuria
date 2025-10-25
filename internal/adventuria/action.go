package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type Action interface {
	User() User
	Type() ActionType
	CanDo() bool
	Do(ActionRequest) (*ActionResult, error)

	setType(ActionType)
	setUser(User)
}

type ActionRecord interface {
	core.RecordProxy

	ID() string
	Save() error
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
	SetNotAffectNextStep(bool)
	DiceRoll() int
	SetDiceRoll(int)
	ItemsUsed() []string
	SetItemsUsed([]string)
	ItemsList() ([]string, error)
	SetItemsList([]string)
	CanMove() bool
	SetCanMove(bool)
}

type ActionType string

const (
	ActionTypeNone ActionType = "none"
	ActionTypeMove ActionType = "move"
)

type ActionRequest struct {
	Comment string
}

type ActionResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
