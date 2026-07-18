package model

import "adventuria/pkg/event"

type Events struct {
	unsubs map[string][]event.Unsubscribe

	onAfterChooseGame      *event.Hook[*OnAfterChooseGameEvent]
	onAfterReroll          *event.Hook[*OnAfterRerollEvent]
	onBeforeDrop           *event.Hook[*OnBeforeDropEvent]
	onBeforeDropCheck      *event.Hook[*OnBeforeDropCheckEvent]
	onAfterDrop            *event.Hook[*OnAfterDropEvent]
	onAfterGoToJail        *event.Hook[*OnAfterGoToJailEvent]
	onBeforeDone           *event.Hook[*OnBeforeDoneEvent]
	onCompleteActivityView *event.Hook[*OnCompleteActivityView]
	onAfterDone            *event.Hook[*OnAfterDoneEvent]
	onBeforeRerollCheck    *event.Hook[*OnBeforeRerollCheckEvent]
	onBeforeRoll           *event.Hook[*OnBeforeRollEvent]
	onBeforeRollMove       *event.Hook[*OnBeforeRollMoveEvent]
	onAfterRoll            *event.Hook[*OnAfterRollEvent]
	onAfterWheelRoll       *event.Hook[*OnAfterWheelRollEvent]
	onAfterItemRoll        *event.Hook[*OnAfterItemRollEvent]
	onAfterItemUse         *event.Hook[*OnAfterItemUseEvent]
	onNewLap               *event.Hook[*OnNewLapEvent]
	onBeforeNextStep       *event.Hook[*OnBeforeNextStepEvent]
	onAfterAction          *event.Hook[*OnAfterActionEvent]
	onAfterMove            *event.Hook[*OnAfterMoveEvent]
	onBeforeCurrentCell    *event.Hook[*OnBeforeCurrentCellEvent]
	onBeforeItemAdd        *event.Hook[*OnBeforeItemAddEvent]
	onAfterItemAdd         *event.Hook[*OnAfterItemAddEvent]
	onBeforeItemBuy        *event.Hook[*OnBeforeItemBuyEvent]
	onBuyGetVariants       *event.Hook[*OnBuyGetViewEvent]
	onBeforeTeleportOnCell *event.Hook[*OnBeforeTeleportOnCellEvent]
	onWorldChanged         *event.Hook[*OnWorldChangedEvent]
}

func NewEvents() *Events {
	return &Events{
		onAfterChooseGame:      &event.Hook[*OnAfterChooseGameEvent]{},
		onAfterReroll:          &event.Hook[*OnAfterRerollEvent]{},
		onBeforeDrop:           &event.Hook[*OnBeforeDropEvent]{},
		onBeforeDropCheck:      &event.Hook[*OnBeforeDropCheckEvent]{},
		onAfterDrop:            &event.Hook[*OnAfterDropEvent]{},
		onAfterGoToJail:        &event.Hook[*OnAfterGoToJailEvent]{},
		onBeforeDone:           &event.Hook[*OnBeforeDoneEvent]{},
		onCompleteActivityView: &event.Hook[*OnCompleteActivityView]{},
		onAfterDone:            &event.Hook[*OnAfterDoneEvent]{},
		onBeforeRerollCheck:    &event.Hook[*OnBeforeRerollCheckEvent]{},
		onBeforeRoll:           &event.Hook[*OnBeforeRollEvent]{},
		onBeforeRollMove:       &event.Hook[*OnBeforeRollMoveEvent]{},
		onAfterRoll:            &event.Hook[*OnAfterRollEvent]{},
		onAfterWheelRoll:       &event.Hook[*OnAfterWheelRollEvent]{},
		onAfterItemRoll:        &event.Hook[*OnAfterItemRollEvent]{},
		onAfterItemUse:         &event.Hook[*OnAfterItemUseEvent]{},
		onNewLap:               &event.Hook[*OnNewLapEvent]{},
		onBeforeNextStep:       &event.Hook[*OnBeforeNextStepEvent]{},
		onAfterAction:          &event.Hook[*OnAfterActionEvent]{},
		onAfterMove:            &event.Hook[*OnAfterMoveEvent]{},
		onBeforeCurrentCell:    &event.Hook[*OnBeforeCurrentCellEvent]{},
		onBeforeItemAdd:        &event.Hook[*OnBeforeItemAddEvent]{},
		onAfterItemAdd:         &event.Hook[*OnAfterItemAddEvent]{},
		onBeforeItemBuy:        &event.Hook[*OnBeforeItemBuyEvent]{},
		onBuyGetVariants:       &event.Hook[*OnBuyGetViewEvent]{},
		onBeforeTeleportOnCell: &event.Hook[*OnBeforeTeleportOnCellEvent]{},
		onWorldChanged:         &event.Hook[*OnWorldChangedEvent]{},
	}
}

func (e *Events) AddUnsubs(key string, unsubs ...event.Unsubscribe) {
	if e.unsubs == nil {
		e.unsubs = make(map[string][]event.Unsubscribe)
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

func (e *Events) OnAfterChooseGame() *event.Hook[*OnAfterChooseGameEvent] {
	return e.onAfterChooseGame
}

func (e *Events) OnAfterReroll() *event.Hook[*OnAfterRerollEvent] {
	return e.onAfterReroll
}

func (e *Events) OnBeforeDrop() *event.Hook[*OnBeforeDropEvent] {
	return e.onBeforeDrop
}

func (e *Events) OnBeforeDropCheck() *event.Hook[*OnBeforeDropCheckEvent] {
	return e.onBeforeDropCheck
}

func (e *Events) OnAfterDrop() *event.Hook[*OnAfterDropEvent] {
	return e.onAfterDrop
}

func (e *Events) OnAfterGoToJail() *event.Hook[*OnAfterGoToJailEvent] {
	return e.onAfterGoToJail
}

func (e *Events) OnBeforeDone() *event.Hook[*OnBeforeDoneEvent] {
	return e.onBeforeDone
}

func (e *Events) OnCompleteActivityView() *event.Hook[*OnCompleteActivityView] {
	return e.onCompleteActivityView
}

func (e *Events) OnAfterDone() *event.Hook[*OnAfterDoneEvent] {
	return e.onAfterDone
}

func (e *Events) OnBeforeRerollCheck() *event.Hook[*OnBeforeRerollCheckEvent] {
	return e.onBeforeRerollCheck
}

func (e *Events) OnBeforeRoll() *event.Hook[*OnBeforeRollEvent] {
	return e.onBeforeRoll
}

func (e *Events) OnBeforeRollMove() *event.Hook[*OnBeforeRollMoveEvent] {
	return e.onBeforeRollMove
}

func (e *Events) OnAfterRoll() *event.Hook[*OnAfterRollEvent] {
	return e.onAfterRoll
}

func (e *Events) OnAfterWheelRoll() *event.Hook[*OnAfterWheelRollEvent] {
	return e.onAfterWheelRoll
}

func (e *Events) OnAfterItemRoll() *event.Hook[*OnAfterItemRollEvent] {
	return e.onAfterItemRoll
}

func (e *Events) OnAfterItemUse() *event.Hook[*OnAfterItemUseEvent] {
	return e.onAfterItemUse
}

func (e *Events) OnNewLap() *event.Hook[*OnNewLapEvent] {
	return e.onNewLap
}

func (e *Events) OnBeforeNextStep() *event.Hook[*OnBeforeNextStepEvent] {
	return e.onBeforeNextStep
}

func (e *Events) OnAfterAction() *event.Hook[*OnAfterActionEvent] {
	return e.onAfterAction
}

func (e *Events) OnAfterMove() *event.Hook[*OnAfterMoveEvent] {
	return e.onAfterMove
}

func (e *Events) OnBeforeCurrentCell() *event.Hook[*OnBeforeCurrentCellEvent] {
	return e.onBeforeCurrentCell
}

func (e *Events) OnBeforeItemAdd() *event.Hook[*OnBeforeItemAddEvent] {
	return e.onBeforeItemAdd
}

func (e *Events) OnAfterItemAdd() *event.Hook[*OnAfterItemAddEvent] {
	return e.onAfterItemAdd
}

func (e *Events) OnBeforeItemBuy() *event.Hook[*OnBeforeItemBuyEvent] {
	return e.onBeforeItemBuy
}

func (e *Events) OnBuyGetView() *event.Hook[*OnBuyGetViewEvent] {
	return e.onBuyGetVariants
}

func (e *Events) OnBeforeTeleportOnCell() *event.Hook[*OnBeforeTeleportOnCellEvent] {
	return e.onBeforeTeleportOnCell
}

func (e *Events) OnWorldChanged() *event.Hook[*OnWorldChangedEvent] {
	return e.onWorldChanged
}

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
type OnBeforeRerollCheckEvent struct {
	event.Event
	IsRerollBlocked bool
}
type OnAfterDropEvent struct {
	event.Event
}
type OnAfterGoToJailEvent struct {
	event.Event
}
type OnBeforeDoneEvent struct {
	event.Event
	CellPoints        int
	CellEnergyConsume int
	CellCoins         int
}
type OnCompleteActivityView struct {
	event.Event
	CellPoints        int
	CellEnergyConsume int
	CellCoins         int
}
type OnAfterDoneEvent struct {
	event.Event
	CurrentCell *CellInfo
}
type OnBeforeRollEvent struct {
	event.Event
	Dices []Dice
}
type OnBeforeRollMoveEvent struct {
	event.Event
	N int
}
type OnAfterRollEvent struct {
	event.Event
	Dices []Dice
	N     int
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
	Data      map[string]any
}
type OnBeforeItemAddEvent struct {
	event.Event
	ItemRecord    *Item
	ShouldAddItem bool
}
type OnAfterItemAddEvent struct {
	event.Event
	Item *InventoryItem
}
type OnNewLapEvent struct {
	event.Event
	Laps int
}
type OnBeforeNextStepEvent struct {
	event.Event
	NextStepType string
	CurrentCell  *CellInfo
}
type OnAfterActionEvent struct {
	event.Event
	ActionType ActionType
}
type OnAfterMoveEvent struct {
	event.Event
	Steps          int
	TotalSteps     int
	PrevTotalSteps int
	CurrentCell    *CellInfo
	CurrentWorld   *World
	Laps           int
}
type OnBeforeCurrentCellEvent struct {
	event.Event
	CurrentCell *CellInfo
}
type OnBeforeItemBuyEvent struct {
	event.Event
	Item  *Item
	Price int
}
type OnBuyGetViewEvent struct {
	event.Event
	Item  *Item
	Price int
}
type OnBeforeTeleportOnCellEvent struct {
	event.Event
	CellId       string
	SkipTeleport bool
}
type OnWorldChangedEvent struct {
	event.Event
	OldWorldId string
	NewWorldId string
}
