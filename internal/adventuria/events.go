package adventuria

import (
	"adventuria/pkg/event"
)

type OnAfterChooseGameEvent struct {
	event.Event
}
type OnAfterRerollEvent struct {
	event.Event
	AppContext AppContext
}
type OnBeforeDropEvent struct {
	event.Event
	AppContext    AppContext
	IsSafeDrop    bool
	IsDropBlocked bool
	PointsForDrop int
}
type OnBeforeDropCheckEvent struct {
	event.Event
	AppContext    AppContext
	IsDropBlocked bool
}
type OnBeforeRerollCheckEvent struct {
	event.Event
	AppContext      AppContext
	IsRerollBlocked bool
}
type OnAfterDropEvent struct {
	event.Event
	AppContext AppContext
}
type OnAfterGoToJailEvent struct {
	event.Event
}
type OnBeforeDoneEvent struct {
	event.Event
	AppContext AppContext
	CellPoints int
	CellCoins  int
}
type OnAfterDoneEvent struct {
	event.Event
	AppContext AppContext
}
type OnBeforeRollEvent struct {
	event.Event
	AppContext AppContext
	Dices      []*Dice
}
type OnBeforeRollMoveEvent struct {
	event.Event
	AppContext AppContext
	N          int
}
type OnAfterRollEvent struct {
	event.Event
	AppContext AppContext
	Dices      []*Dice
	N          int
}
type OnBeforeWheelRollEvent struct {
	event.Event
	AppContext  AppContext
	CurrentCell CellWheel
}
type OnAfterWheelRollEvent struct {
	event.Event
	AppContext AppContext
	ItemId     string
}
type OnAfterItemRollEvent struct {
	event.Event
	AppContext AppContext
	ItemId     string
}
type OnAfterItemUseEvent struct {
	event.Event
	AppContext AppContext
	InvItemId  string
	Data       map[string]any
}
type OnBeforeItemAdd struct {
	event.Event
	AppContext    AppContext
	ItemRecord    ItemRecord
	ShouldAddItem bool
}
type OnAfterItemAdd struct {
	event.Event
	AppContext AppContext
	ItemRecord ItemRecord
}
type OnAfterItemSave struct {
	event.Event
	AppContext AppContext
	Item       Item
}
type OnNewLapEvent struct {
	event.Event
	AppContext AppContext
	Laps       int
}
type OnBeforeNextStepEvent struct {
	event.Event
	NextStepType string
	CurrentCell  Cell
}
type OnAfterActionEvent struct {
	event.Event
	AppContext AppContext
	ActionType ActionType
}
type OnAfterMoveEvent struct {
	event.Event
	AppContext     AppContext
	Steps          int
	TotalSteps     int
	PrevTotalSteps int
	CurrentCell    Cell
	Laps           int
}
type OnBeforeCurrentCellEvent struct {
	event.Event
	CurrentCell Cell
}
type OnBeforeItemBuy struct {
	event.Event
	AppContext AppContext
	Item       ItemRecord
	Price      int
}
type OnBuyGetVariants struct {
	event.Event
	AppContext AppContext
	Item       ItemRecord
	Price      int
}
