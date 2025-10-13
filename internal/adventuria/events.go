package adventuria

import (
	"adventuria/pkg/event"
)

type OnAfterChooseGameEvent struct {
	event.Event
}
type OnAfterRerollEvent struct {
	event.Event
}

type OnBeforeDropEvent struct {
	event.Event
	IsSafeDrop bool
}
type OnAfterDropEvent struct {
	event.Event
}
type OnAfterGoToJailEvent struct {
	event.Event
}
type OnBeforeDoneEvent struct {
	event.Event
	CellPointsDivide int
}
type OnAfterDoneEvent struct {
	event.Event
}
type OnBeforeRollEvent struct {
	event.Event
	Dices []*Dice
}
type OnBeforeRollMoveEvent struct {
	event.Event
	N int
}
type OnAfterRollEvent struct {
	event.Event
	Dices []*Dice
	N     int
}
type OnBeforeWheelRollEvent struct {
	event.Event
	CurrentCell CellWheel
}
type OnAfterWheelRollEvent struct {
	event.Event
	ItemId string
}
type OnAfterItemRollEvent struct {
	event.Event
	ItemId string
}
type OnAfterItemUseEvent struct {
	event.Event
	ItemId string
}
type OnNewLapEvent struct {
	event.Event
	Laps int
}
type OnBeforeNextStepEvent struct {
	event.Event
	NextStepType string
	CurrentCell  Cell
}
type OnAfterActionEvent struct {
	event.Event
}
type OnAfterMoveEvent struct {
	event.Event
	Steps       int
	CurrentCell Cell
	Laps        int
}
type OnBeforeCurrentCellEvent struct {
	event.Event
	CurrentCell Cell
}
