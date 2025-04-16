package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	ActionTypeRollDice     = "rollDice"
	ActionTypeDone         = "done"
	ActionTypeReroll       = "reroll"
	ActionTypeDrop         = "drop"
	ActionTypeChooseResult = "chooseResult"
	ActionTypeRollWheel    = "rollWheel"
)

type Action interface {
	core.RecordProxy
	Save() error
	UserId() string
	CellId() string
	SetCell(string)
	Comment() string
	SetComment(string)
	Value() string
	SetValue(value any)
	Type() string
	SetType(string)
	SetNotAffectNextStep(bool)
	CollectionRef() string
	SetCollectionRef(string)
	DiceRoll() int
	SetDiceRoll(int)
	ItemsUsed() []string
	SetItemsUsed([]string)
}
