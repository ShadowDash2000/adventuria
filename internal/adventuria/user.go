package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/event"

	"github.com/pocketbase/pocketbase/core"
)

type User interface {
	core.RecordProxy
	UserEvent
	cache.Closable

	ID() string
	Name() string
	IsSafeDrop() bool
	IsInJail() bool
	SetIsInJail(b bool)
	CurrentCell() (Cell, bool)
	Points() int
	SetPoints(points int)
	DropsInARow() int
	SetDropsInARow(drops int)
	CellsPassed() int
	setCellsPassed(cellsPassed int)
	MaxInventorySlots() int
	SetMaxInventorySlots(maxInventorySlots int)
	ItemWheelsCount() int
	SetItemWheelsCount(itemWheelsCount int)
	Move(steps int) (*OnAfterMoveEvent, error)
	MoveToCellType(cellType CellType) error
	MoveToCellId(cellId string) error
	MoveToCellName(cellName string) error
	Inventory() Inventory
	LastAction() ActionRecord
	Timer() Timer
	Stats() *Stats
	Balance() int
	SetBalance(int)
	IsStreamLive() bool
	SetIsStreamLive(bool)
	isInAction() bool
	setIsInAction(bool)
}

type UserEvent interface {
	OnAfterChooseGame() *event.Hook[*OnAfterChooseGameEvent]
	OnAfterReroll() *event.Hook[*OnAfterRerollEvent]
	OnBeforeDrop() *event.Hook[*OnBeforeDropEvent]
	OnAfterDrop() *event.Hook[*OnAfterDropEvent]
	OnAfterGoToJail() *event.Hook[*OnAfterGoToJailEvent]
	OnBeforeDone() *event.Hook[*OnBeforeDoneEvent]
	OnAfterDone() *event.Hook[*OnAfterDoneEvent]
	OnBeforeRoll() *event.Hook[*OnBeforeRollEvent]
	OnBeforeRollMove() *event.Hook[*OnBeforeRollMoveEvent]
	OnAfterRoll() *event.Hook[*OnAfterRollEvent]
	OnBeforeWheelRoll() *event.Hook[*OnBeforeWheelRollEvent]
	OnAfterWheelRoll() *event.Hook[*OnAfterWheelRollEvent]
	OnAfterItemRoll() *event.Hook[*OnAfterItemRollEvent]
	OnAfterItemUse() *event.Hook[*OnAfterItemUseEvent]
	OnNewLap() *event.Hook[*OnNewLapEvent]
	OnBeforeNextStep() *event.Hook[*OnBeforeNextStepEvent]
	OnAfterAction() *event.Hook[*OnAfterActionEvent]
	OnAfterMove() *event.Hook[*OnAfterMoveEvent]
	OnBeforeCurrentCell() *event.Hook[*OnBeforeCurrentCellEvent]
}

type Stats struct {
	Drops       int `json:"drops"`
	Rerolls     int `json:"rerolls"`
	Finished    int `json:"finished"`
	WasInJail   int `json:"wasInJail"`
	ItemsUsed   int `json:"itemsUsed"`
	DiceRolls   int `json:"diceRolls"`
	MaxDiceRoll int `json:"maxDiceRoll"`
	WheelRolled int `json:"wheelRolled"`
}
