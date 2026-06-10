package model

import (
	"adventuria/pkg/event_new"
)

type Events struct {
	unsubs map[string][]event_new.Unsubscribe

	onAfterChooseGame      *event_new.Hook[*OnAfterChooseGameEvent]
	onAfterReroll          *event_new.Hook[*OnAfterRerollEvent]
	onBeforeDrop           *event_new.Hook[*OnBeforeDropEvent]
	onBeforeDropCheck      *event_new.Hook[*OnBeforeDropCheckEvent]
	onAfterDrop            *event_new.Hook[*OnAfterDropEvent]
	onAfterGoToJail        *event_new.Hook[*OnAfterGoToJailEvent]
	onBeforeDone           *event_new.Hook[*OnBeforeDoneEvent]
	onAfterDone            *event_new.Hook[*OnAfterDoneEvent]
	onBeforeRerollCheck    *event_new.Hook[*OnBeforeRerollCheckEvent]
	onBeforeRoll           *event_new.Hook[*OnBeforeRollEvent]
	onBeforeRollMove       *event_new.Hook[*OnBeforeRollMoveEvent]
	onAfterRoll            *event_new.Hook[*OnAfterRollEvent]
	onAfterWheelRoll       *event_new.Hook[*OnAfterWheelRollEvent]
	onAfterItemRoll        *event_new.Hook[*OnAfterItemRollEvent]
	onAfterItemUse         *event_new.Hook[*OnAfterItemUseEvent]
	onNewLap               *event_new.Hook[*OnNewLapEvent]
	onBeforeNextStep       *event_new.Hook[*OnBeforeNextStepEvent]
	onAfterAction          *event_new.Hook[*OnAfterActionEvent]
	onAfterMove            *event_new.Hook[*OnAfterMoveEvent]
	onBeforeCurrentCell    *event_new.Hook[*OnBeforeCurrentCellEvent]
	onBeforeItemAdd        *event_new.Hook[*OnBeforeItemAddEvent]
	onAfterItemAdd         *event_new.Hook[*OnAfterItemAddEvent]
	onBeforeItemBuy        *event_new.Hook[*OnBeforeItemBuyEvent]
	onBuyGetVariants       *event_new.Hook[*OnBuyGetViewEvent]
	onBeforeTeleportOnCell *event_new.Hook[*OnBeforeTeleportOnCellEvent]
	onWorldChanged         *event_new.Hook[*OnWorldChangedEvent]
}

func NewEvents() *Events {
	return &Events{
		onAfterChooseGame:      &event_new.Hook[*OnAfterChooseGameEvent]{},
		onAfterReroll:          &event_new.Hook[*OnAfterRerollEvent]{},
		onBeforeDrop:           &event_new.Hook[*OnBeforeDropEvent]{},
		onBeforeDropCheck:      &event_new.Hook[*OnBeforeDropCheckEvent]{},
		onAfterDrop:            &event_new.Hook[*OnAfterDropEvent]{},
		onAfterGoToJail:        &event_new.Hook[*OnAfterGoToJailEvent]{},
		onBeforeDone:           &event_new.Hook[*OnBeforeDoneEvent]{},
		onAfterDone:            &event_new.Hook[*OnAfterDoneEvent]{},
		onBeforeRerollCheck:    &event_new.Hook[*OnBeforeRerollCheckEvent]{},
		onBeforeRoll:           &event_new.Hook[*OnBeforeRollEvent]{},
		onBeforeRollMove:       &event_new.Hook[*OnBeforeRollMoveEvent]{},
		onAfterRoll:            &event_new.Hook[*OnAfterRollEvent]{},
		onAfterWheelRoll:       &event_new.Hook[*OnAfterWheelRollEvent]{},
		onAfterItemRoll:        &event_new.Hook[*OnAfterItemRollEvent]{},
		onAfterItemUse:         &event_new.Hook[*OnAfterItemUseEvent]{},
		onNewLap:               &event_new.Hook[*OnNewLapEvent]{},
		onBeforeNextStep:       &event_new.Hook[*OnBeforeNextStepEvent]{},
		onAfterAction:          &event_new.Hook[*OnAfterActionEvent]{},
		onAfterMove:            &event_new.Hook[*OnAfterMoveEvent]{},
		onBeforeCurrentCell:    &event_new.Hook[*OnBeforeCurrentCellEvent]{},
		onBeforeItemAdd:        &event_new.Hook[*OnBeforeItemAddEvent]{},
		onAfterItemAdd:         &event_new.Hook[*OnAfterItemAddEvent]{},
		onBeforeItemBuy:        &event_new.Hook[*OnBeforeItemBuyEvent]{},
		onBuyGetVariants:       &event_new.Hook[*OnBuyGetViewEvent]{},
		onBeforeTeleportOnCell: &event_new.Hook[*OnBeforeTeleportOnCellEvent]{},
		onWorldChanged:         &event_new.Hook[*OnWorldChangedEvent]{},
	}
}

func (e *Events) AddUnsubs(key string, unsubs ...event_new.Unsubscribe) {
	if e.unsubs == nil {
		e.unsubs = make(map[string][]event_new.Unsubscribe)
	}
	e.unsubs[key] = append(e.unsubs[key], unsubs...)
}

func (e *Events) Unsubscribe(key string) {
	if unsubs, ok := e.unsubs[key]; ok {
		for _, unsub := range unsubs {
			unsub()
		}
		delete(e.unsubs, key)
	}
}

func (e *Events) Close() {
	for _, unsubs := range e.unsubs {
		for _, unsub := range unsubs {
			unsub()
		}
	}
	e.unsubs = nil
}

func (e *Events) OnAfterChooseGame() *event_new.Hook[*OnAfterChooseGameEvent] {
	return e.onAfterChooseGame
}

func (e *Events) OnAfterReroll() *event_new.Hook[*OnAfterRerollEvent] {
	return e.onAfterReroll
}

func (e *Events) OnBeforeDrop() *event_new.Hook[*OnBeforeDropEvent] {
	return e.onBeforeDrop
}

func (e *Events) OnBeforeDropCheck() *event_new.Hook[*OnBeforeDropCheckEvent] {
	return e.onBeforeDropCheck
}

func (e *Events) OnAfterDrop() *event_new.Hook[*OnAfterDropEvent] {
	return e.onAfterDrop
}

func (e *Events) OnAfterGoToJail() *event_new.Hook[*OnAfterGoToJailEvent] {
	return e.onAfterGoToJail
}

func (e *Events) OnBeforeDone() *event_new.Hook[*OnBeforeDoneEvent] {
	return e.onBeforeDone
}

func (e *Events) OnAfterDone() *event_new.Hook[*OnAfterDoneEvent] {
	return e.onAfterDone
}

func (e *Events) OnBeforeRerollCheck() *event_new.Hook[*OnBeforeRerollCheckEvent] {
	return e.onBeforeRerollCheck
}

func (e *Events) OnBeforeRoll() *event_new.Hook[*OnBeforeRollEvent] {
	return e.onBeforeRoll
}

func (e *Events) OnBeforeRollMove() *event_new.Hook[*OnBeforeRollMoveEvent] {
	return e.onBeforeRollMove
}

func (e *Events) OnAfterRoll() *event_new.Hook[*OnAfterRollEvent] {
	return e.onAfterRoll
}

func (e *Events) OnAfterWheelRoll() *event_new.Hook[*OnAfterWheelRollEvent] {
	return e.onAfterWheelRoll
}

func (e *Events) OnAfterItemRoll() *event_new.Hook[*OnAfterItemRollEvent] {
	return e.onAfterItemRoll
}

func (e *Events) OnAfterItemUse() *event_new.Hook[*OnAfterItemUseEvent] {
	return e.onAfterItemUse
}

func (e *Events) OnNewLap() *event_new.Hook[*OnNewLapEvent] {
	return e.onNewLap
}

func (e *Events) OnBeforeNextStep() *event_new.Hook[*OnBeforeNextStepEvent] {
	return e.onBeforeNextStep
}

func (e *Events) OnAfterAction() *event_new.Hook[*OnAfterActionEvent] {
	return e.onAfterAction
}

func (e *Events) OnAfterMove() *event_new.Hook[*OnAfterMoveEvent] {
	return e.onAfterMove
}

func (e *Events) OnBeforeCurrentCell() *event_new.Hook[*OnBeforeCurrentCellEvent] {
	return e.onBeforeCurrentCell
}

func (e *Events) OnBeforeItemAdd() *event_new.Hook[*OnBeforeItemAddEvent] {
	return e.onBeforeItemAdd
}

func (e *Events) OnAfterItemAdd() *event_new.Hook[*OnAfterItemAddEvent] {
	return e.onAfterItemAdd
}

func (e *Events) OnBeforeItemBuy() *event_new.Hook[*OnBeforeItemBuyEvent] {
	return e.onBeforeItemBuy
}

func (e *Events) OnBuyGetView() *event_new.Hook[*OnBuyGetViewEvent] {
	return e.onBuyGetVariants
}

func (e *Events) OnBeforeTeleportOnCell() *event_new.Hook[*OnBeforeTeleportOnCellEvent] {
	return e.onBeforeTeleportOnCell
}

func (e *Events) OnWorldChanged() *event_new.Hook[*OnWorldChangedEvent] {
	return e.onWorldChanged
}

type OnAfterChooseGameEvent struct {
	event_new.Event
}
type OnAfterRerollEvent struct {
	event_new.Event
}
type OnBeforeDropEvent struct {
	event_new.Event
	IsSafeDrop    bool
	IsDropBlocked bool
	PointsForDrop int
}
type OnBeforeDropCheckEvent struct {
	event_new.Event
	IsDropBlocked bool
}
type OnBeforeRerollCheckEvent struct {
	event_new.Event
	IsRerollBlocked bool
}
type OnAfterDropEvent struct {
	event_new.Event
}
type OnAfterGoToJailEvent struct {
	event_new.Event
}
type OnBeforeDoneEvent struct {
	event_new.Event
	CellPoints int
	CellCoins  int
}
type OnAfterDoneEvent struct {
	event_new.Event
}
type OnBeforeRollEvent struct {
	event_new.Event
	Dices []Dice
}
type OnBeforeRollMoveEvent struct {
	event_new.Event
	N int
}
type OnAfterRollEvent struct {
	event_new.Event
	Dices []Dice
	N     int
}
type OnAfterWheelRollEvent struct {
	event_new.Event
	ItemId string
}
type OnAfterItemRollEvent struct {
	event_new.Event
	ItemId string
}
type OnAfterItemUseEvent struct {
	event_new.Event
	InvItemId string
	Data      map[string]any
}
type OnBeforeItemAddEvent struct {
	event_new.Event
	ItemRecord    *Item
	ShouldAddItem bool
}
type OnAfterItemAddEvent struct {
	event_new.Event
	Item *InventoryItem
}
type OnNewLapEvent struct {
	event_new.Event
	Laps int
}
type OnBeforeNextStepEvent struct {
	event_new.Event
	NextStepType string
	CurrentCell  *CellInfo
}
type OnAfterActionEvent struct {
	event_new.Event
	ActionType ActionType
}
type OnAfterMoveEvent struct {
	event_new.Event
	Steps          int
	TotalSteps     int
	PrevTotalSteps int
	CurrentCell    *CellInfo
	CellLocalOrder int
	CurrentWorld   *World
	Laps           int
}
type OnBeforeCurrentCellEvent struct {
	event_new.Event
	CurrentCell *CellInfo
}
type OnBeforeItemBuyEvent struct {
	event_new.Event
	Item  *Item
	Price int
}
type OnBuyGetViewEvent struct {
	event_new.Event
	Item  *Item
	Price int
}
type OnBeforeTeleportOnCellEvent struct {
	event_new.Event
	CellId       string
	SkipTeleport bool
}
type OnWorldChangedEvent struct {
	event_new.Event
	OldWorldId string
	NewWorldId string
}
