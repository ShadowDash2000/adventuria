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
	IsSafeDrop    bool
	IsDropBlocked bool
	PointsForDrop int
}
type OnBeforeDropCheckEvent struct {
	event.Event
	IsDropBlocked bool
}
type OnAfterDropEvent struct {
	event.Event
}
type OnAfterGoToJailEvent struct {
	event.Event
}
type OnBeforeDoneEvent struct {
	event.Event
	CellPoints int
	CellCoins  int
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
	InvItemId string
	Request   UseItemRequest
}
type OnBeforeItemAdd struct {
	event.Event
	ItemRecord ItemRecord
}
type OnAfterItemAdd struct {
	event.Event
	ItemRecord ItemRecord
}
type OnAfterItemSave struct {
	event.Event
	Item Item
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
	ActionType ActionType
}
type OnAfterMoveEvent struct {
	event.Event
	Steps          int  `json:"steps"`
	TotalSteps     int  `json:"total_steps"`
	PrevTotalSteps int  `json:"prev_total_steps"`
	CurrentCell    Cell `json:"current_cell"`
	Laps           int  `json:"laps"`
}
type OnBeforeCurrentCellEvent struct {
	event.Event
	CurrentCell Cell
}
