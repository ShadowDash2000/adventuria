package adventuria

import (
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type Action interface {
	core.RecordProxy
	Save() error
	User() User
	UserId() string
	CellId() string
	SetCell(string)
	Comment() string
	SetComment(string)
	Value() string
	SetValue(value any)
	Type() ActionType
	SetType(ActionType)
	SetNotAffectNextStep(bool)
	CollectionRef() string
	SetCollectionRef(string)
	DiceRoll() int
	SetDiceRoll(int)
	ItemsUsed() []string
	SetItemsUsed([]string)

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

var actionsList = map[ActionType]ActionCreator{}

type ActionCreator func() Action

func RegisterActions(actions map[ActionType]ActionCreator) {
	maps.Insert(actionsList, maps.All(actions))
}

func NewAction(a Action) ActionCreator {
	return func() Action {
		return a
	}
}
