package adventuria

import (
	"github.com/AlexanderGrom/go-event"
)

const (
	OnAfterChooseGame    = "OnAfterChooseGame"
	OnAfterReroll        = "OnAfterReroll"
	OnBeforeDrop         = "OnBeforeDrop"
	OnAfterDrop          = "OnAfterDrop"
	OnAfterGoToJail      = "OnAfterGoToJail"
	OnBeforeDone         = "OnBeforeDone"
	OnAfterDone          = "OnAfterDone"
	OnBeforeRoll         = "OnBeforeRoll"
	OnBeforeRollMove     = "OnBeforeRollMove"
	OnAfterRoll          = "OnAfterRoll"
	OnBeforeWheelRoll    = "OnBeforeWheelRoll"
	OnAfterWheelRoll     = "OnAfterWheelRoll"
	OnAfterItemRoll      = "OnAfterItemRoll"
	OnAfterItemUse       = "OnAfterItemUse"
	OnNewLap             = "OnNewLap"
	OnBeforeNextStepType = "OnBeforeNextStepType"
	OnAfterAction        = "OnAfterAction"
	OnAfterMove          = "OnAfterMove"
)

type Event interface {
	On(string, EventFn) error
	Go(string, EventFields) error
}

type EventBase struct {
	e event.Dispatcher
}

func NewEvent() Event {
	return &EventBase{
		e: event.New(),
	}
}

type EventFn func(EventFields) error

func (e *EventBase) On(name string, fn EventFn) error {
	return e.e.On(name, fn)
}

func (e *EventBase) Go(name string, fields EventFields) error {
	return e.e.Go(name, fields)
}

type EventFields interface {
	User() *User
	Fields() any
	Effects(EffectUse) *Effects
}

type EventFieldsBase struct {
	user    *User
	fields  any
	effects map[EffectUse]*Effects
}

func NewEventFields(user *User, fields any) EventFields {
	return &EventFieldsBase{
		user:    user,
		fields:  fields,
		effects: map[EffectUse]*Effects{},
	}
}

func (e *EventFieldsBase) User() *User {
	return e.user
}

func (e *EventFieldsBase) Fields() any {
	return e.fields
}

func (e *EventFieldsBase) Effects(event EffectUse) *Effects {
	if effects, ok := e.effects[event]; ok {
		return effects
	}

	e.effects[event], _, _ = e.User().Inventory.Effects(event)
	return e.effects[event]
}

func (e *EventFieldsBase) Effect(event EffectUse, effect string) Effect {
	effects := e.Effects(event)
	return effects.Effect(effect)
}

type OnAfterChooseGameFields struct {
}
type OnAfterRerollFields struct {
}
type OnBeforeDropFields struct {
	IsSafeDrop bool
}
type OnAfterDropFields struct {
}
type OnAfterGoToJailFields struct {
}
type OnBeforeDoneFields struct {
	CellPointsDivide int
}
type OnAfterDoneFields struct {
}
type OnBeforeRollFields struct {
	Dices []Dice
}
type OnBeforeRollMoveFields struct {
	N int
}
type OnAfterRollFields struct {
	Dices []Dice
	N     int
}
type OnBeforeWheelRollFields struct {
	CurrentCell CellWheel
}
type OnAfterWheelRollFields struct {
	ItemId string
}
type OnAfterItemRollFields struct {
	ItemId string
}
type OnAfterItemUseFields struct {
	ItemId string
}
type OnNewLapFields struct {
	Laps int
}
type OnBeforeNextStepFields struct {
	NextStepType string
	CurrentCell  Cell
}
type OnAfterActionFields struct {
	Event EffectUse
}
type OnAfterMoveFields struct {
	Steps       int
	Action      Action
	CurrentCell Cell
	Laps        int
}
