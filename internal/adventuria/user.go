package adventuria

import (
	"adventuria/pkg/event"

	"github.com/pocketbase/pocketbase/core"
)

type User interface {
	core.RecordProxy
	UserEvent
	Closable

	Refetch(ctx AppContext) error
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
	addCellsPassed(ctx AppContext, amount int) error
	MaxInventorySlots() int
	SetMaxInventorySlots(maxInventorySlots int)
	ItemWheelsCount() int
	SetItemWheelsCount(itemWheelsCount int)
	Move(ctx AppContext, steps int) ([]*MoveResult, error)
	CurrentCellOrder() int
	MoveToClosestCellType(ctx AppContext, cellType CellType) ([]*MoveResult, error)
	MoveToCellId(ctx AppContext, cellId string) ([]*MoveResult, error)
	MoveToCellName(ctx AppContext, cellName string) ([]*MoveResult, error)
	MoveToClosestCellByNames(ctx AppContext, cellNames ...string) ([]*MoveResult, error)
	Inventory() Inventory
	LastAction() ActionRecord
	Timer() Timer
	Stats() *Stats
	Balance() int
	AddBalance(AppContext, int) error
	IsStreamLive() bool
	SetIsStreamLive(bool)
	isInAction() bool
	setIsInAction(bool)
}

type UserEvent interface {
	OnAfterChooseGame() *event.Hook[*OnAfterChooseGameEvent]
	OnAfterReroll() *event.Hook[*OnAfterRerollEvent]
	OnBeforeDrop() *event.Hook[*OnBeforeDropEvent]
	OnBeforeDropCheck() *event.Hook[*OnBeforeDropCheckEvent]
	OnAfterDrop() *event.Hook[*OnAfterDropEvent]
	OnAfterGoToJail() *event.Hook[*OnAfterGoToJailEvent]
	OnBeforeDone() *event.Hook[*OnBeforeDoneEvent]
	OnAfterDone() *event.Hook[*OnAfterDoneEvent]
	OnBeforeRerollCheck() *event.Hook[*OnBeforeRerollCheckEvent]
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
	OnBeforeItemAdd() *event.Hook[*OnBeforeItemAdd]
	OnAfterItemAdd() *event.Hook[*OnAfterItemAdd]
	OnAfterItemSave() *event.Hook[*OnAfterItemSave]
	OnBeforeItemBuy() *event.Hook[*OnBeforeItemBuy]
	OnBuyGetVariants() *event.Hook[*OnBuyGetVariants]
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

type MoveResult struct {
	Steps          int  `json:"steps"`
	TotalSteps     int  `json:"total_steps"`
	PrevTotalSteps int  `json:"prev_total_steps"`
	CurrentCell    Cell `json:"current_cell"`
	Laps           int  `json:"laps"`
}
