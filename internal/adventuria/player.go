package adventuria

import (
	"adventuria/pkg/event"

	"github.com/pocketbase/pocketbase/core"
)

type Player interface {
	core.RecordProxy
	PlayerEvent
	Closable

	Refetch(ctx AppContext) error
	ID() string
	Name() string
	IsStreamLive() bool
	SetIsStreamLive(bool)

	Move(ctx AppContext, steps int) ([]*MoveResult, error)
	MoveToClosestCellType(ctx AppContext, cellType CellType) ([]*MoveResult, error)
	MoveToCellId(ctx AppContext, cellId string) ([]*MoveResult, error)

	Inventory() Inventory
	LastAction() ActionRecord
	Progress() PlayerProgress

	Locked() bool
	Lock()
	Unlock()
}

type PlayerEvent interface {
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
	OnBeforeTeleportOnCell() *event.Hook[*OnBeforeTeleportOnCell]
	OnWorldChanged() *event.Hook[*OnWorldChangedEvent]
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
