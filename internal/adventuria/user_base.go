package adventuria

import (
	"adventuria/pkg/event"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

type UserBase struct {
	core.BaseRecordProxy
	locator    ServiceLocator
	lastAction Action
	inventory  Inventory
	timer      Timer
	stats      *Stats

	onAfterChooseGame   *event.Hook[*OnAfterChooseGameEvent]
	onAfterReroll       *event.Hook[*OnAfterRerollEvent]
	onBeforeDrop        *event.Hook[*OnBeforeDropEvent]
	onAfterDrop         *event.Hook[*OnAfterDropEvent]
	onAfterGoToJail     *event.Hook[*OnAfterGoToJailEvent]
	onBeforeDone        *event.Hook[*OnBeforeDoneEvent]
	onAfterDone         *event.Hook[*OnAfterDoneEvent]
	onBeforeRoll        *event.Hook[*OnBeforeRollEvent]
	onBeforeRollMove    *event.Hook[*OnBeforeRollMoveEvent]
	onAfterRoll         *event.Hook[*OnAfterRollEvent]
	onBeforeWheelRoll   *event.Hook[*OnBeforeWheelRollEvent]
	onAfterWheelRoll    *event.Hook[*OnAfterWheelRollEvent]
	onAfterItemRoll     *event.Hook[*OnAfterItemRollEvent]
	onAfterItemUse      *event.Hook[*OnAfterItemUseEvent]
	onNewLap            *event.Hook[*OnNewLapEvent]
	onBeforeNextStep    *event.Hook[*OnBeforeNextStepEvent]
	onAfterAction       *event.Hook[*OnAfterActionEvent]
	onAfterMove         *event.Hook[*OnAfterMoveEvent]
	onBeforeCurrentCell *event.Hook[*OnBeforeCurrentCellEvent]
}

func NewUser(locator ServiceLocator, userId string) (User, error) {
	if userId == "" {
		return nil, errors.New("empty user id")
	}

	var err error
	timer, err := NewTimer(locator, userId)
	if err != nil {
		return nil, err
	}

	u := &UserBase{
		locator: locator,
		timer:   timer,
	}

	err = u.fetchUser(userId)
	if err != nil {
		return nil, err
	}

	u.lastAction, err = NewLastUserAction(u.locator, userId)
	if err != nil {
		return nil, err
	}

	u.inventory, err = NewInventory(u.locator, u, u.MaxInventorySlots())
	if err != nil {
		return nil, err
	}

	u.bindHooks()

	return u, nil
}

func (u *UserBase) bindHooks() {
	u.locator.PocketBase().OnRecordAfterUpdateSuccess(TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			u.SetProxyRecord(e.Record)
		}
		return e.Next()
	})

	var persistentEffects []PersistentEffect
	for _, effectCreator := range persistentEffectsList {
		effect := effectCreator(u)
		effect.Subscribe()
		persistentEffects = append(persistentEffects, effect)
	}
}

func (u *UserBase) SetProxyRecord(record *core.Record) {
	u.BaseRecordProxy.SetProxyRecord(record)
	u.inventory.SetMaxSlots(u.MaxInventorySlots())
	u.UnmarshalJSONField("stats", &u.stats)
}

func (u *UserBase) fetchUser(userId string) error {
	user, err := u.locator.PocketBase().FindRecordById(TableUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)

	return nil
}

func (u *UserBase) ID() string {
	return u.Id
}

func (u *UserBase) SetCantDrop(b bool) {
	u.Set("cantDrop", b)
}

func (u *UserBase) CanDrop() bool {
	return !u.GetBool("cantDrop") && !u.IsInJail()
}

func (u *UserBase) IsSafeDrop() bool {
	return u.DropsInARow() < u.locator.Settings().DropsToJail()
}

func (u *UserBase) IsInJail() bool {
	return u.GetBool("isInJail")
}

func (u *UserBase) SetIsInJail(b bool) {
	u.Set("isInJail", b)
}

func (u *UserBase) CurrentCell() (Cell, bool) {
	currentCellNum := u.CellsPassed() % u.locator.Cells().Count()
	cell, ok := u.locator.Cells().GetByOrder(currentCellNum)

	return cell, ok
}

func (u *UserBase) Points() int {
	return u.GetInt("points")
}

func (u *UserBase) SetPoints(points int) {
	u.Set("points", points)
}

func (u *UserBase) DropsInARow() int {
	return u.GetInt("dropsInARow")
}

func (u *UserBase) SetDropsInARow(drops int) {
	u.Set("dropsInARow", drops)
}

func (u *UserBase) CellsPassed() int {
	return u.GetInt("cellsPassed")
}

func (u *UserBase) SetCellsPassed(cellsPassed int) {
	u.Set("cellsPassed", cellsPassed)
}

func (u *UserBase) MaxInventorySlots() int {
	return u.GetInt("maxInventorySlots")
}

func (u *UserBase) SetMaxInventorySlots(maxInventorySlots int) {
	u.Set("maxInventorySlots", maxInventorySlots)
}

func (u *UserBase) ItemWheelsCount() int {
	return u.GetInt("itemWheelsCount")
}

func (u *UserBase) SetItemWheelsCount(itemWheelsCount int) {
	u.Set("itemWheelsCount", itemWheelsCount)
}

func (u *UserBase) save() error {
	u.Set("stats", u.stats)

	return u.locator.PocketBase().Save(u)
}

func (u *UserBase) Move(steps int) (*OnAfterMoveEvent, error) {
	cellsPassed := u.CellsPassed()
	currentCellNum := (cellsPassed + steps) % u.locator.Cells().Count()
	currentCell, _ := u.locator.Cells().GetByOrder(currentCellNum)

	u.SetCellsPassed(cellsPassed + steps)

	err := currentCell.OnCellReached(u)
	if err != nil {
		return nil, err
	}

	action, err := NewActionFromType(u.locator, u, ActionTypeRollDice)
	if err != nil {
		return nil, err
	}
	action.SetCell(currentCell.ID())
	action.SetValue(steps)
	err = action.Save()
	if err != nil {
		return nil, err
	}

	prevCellNum := cellsPassed % u.locator.Cells().Count()
	lapsPassed := (prevCellNum + steps) / u.locator.Cells().Count()
	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		err = u.OnNewLap().Trigger(&OnNewLapEvent{
			Laps: lapsPassed,
		})
		if err != nil {
			return nil, err
		}
	}

	onAfterMoveEvent := &OnAfterMoveEvent{
		Steps:       steps,
		Action:      action,
		CurrentCell: currentCell,
		Laps:        lapsPassed,
	}

	err = u.OnAfterMove().Trigger(onAfterMoveEvent)
	if err != nil {
		return nil, err
	}

	return onAfterMoveEvent, nil
}

func (u *UserBase) MoveToCellType(cellType CellType) error {
	cellPos, ok := u.locator.Cells().GetOrderByType(cellType)
	if !ok {
		return errors.New("cell not found")
	}

	currentCellNum := u.CellsPassed() % u.locator.Cells().Count()

	_, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserBase) MoveToCellId(cellId string) error {
	cellPos, ok := u.locator.Cells().GetOrderById(cellId)
	if !ok {
		return fmt.Errorf("cell %s not found", cellId)
	}

	currentCellNum := u.CellsPassed() % u.locator.Cells().Count()

	_, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

// GetNextStepType
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *UserBase) GetNextStepType() string {
	currentCell, ok := u.CurrentCell()
	if !ok {
		return ""
	}

	return currentCell.NextStep(u)
}

func (u *UserBase) Inventory() Inventory {
	return u.inventory
}

func (u *UserBase) LastAction() Action {
	return u.lastAction
}

func (u *UserBase) Timer() Timer {
	return u.timer
}

func (u *UserBase) Stats() *Stats {
	return u.stats
}

func (u *UserBase) OnAfterChooseGame() *event.Hook[*OnAfterChooseGameEvent] {
	return u.onAfterChooseGame
}

func (u *UserBase) OnAfterReroll() *event.Hook[*OnAfterRerollEvent] {
	return u.onAfterReroll
}

func (u *UserBase) OnBeforeDrop() *event.Hook[*OnBeforeDropEvent] {
	return u.onBeforeDrop
}

func (u *UserBase) OnAfterDrop() *event.Hook[*OnAfterDropEvent] {
	return u.onAfterDrop
}

func (u *UserBase) OnAfterGoToJail() *event.Hook[*OnAfterGoToJailEvent] {
	return u.onAfterGoToJail
}

func (u *UserBase) OnBeforeDone() *event.Hook[*OnBeforeDoneEvent] {
	return u.onBeforeDone
}

func (u *UserBase) OnAfterDone() *event.Hook[*OnAfterDoneEvent] {
	return u.onAfterDone
}

func (u *UserBase) OnBeforeRoll() *event.Hook[*OnBeforeRollEvent] {
	return u.onBeforeRoll
}

func (u *UserBase) OnBeforeRollMove() *event.Hook[*OnBeforeRollMoveEvent] {
	return u.onBeforeRollMove
}

func (u *UserBase) OnAfterRoll() *event.Hook[*OnAfterRollEvent] {
	return u.onAfterRoll
}

func (u *UserBase) OnBeforeWheelRoll() *event.Hook[*OnBeforeWheelRollEvent] {
	return u.onBeforeWheelRoll
}

func (u *UserBase) OnAfterWheelRoll() *event.Hook[*OnAfterWheelRollEvent] {
	return u.onAfterWheelRoll
}

func (u *UserBase) OnAfterItemRoll() *event.Hook[*OnAfterItemRollEvent] {
	return u.onAfterItemRoll
}

func (u *UserBase) OnAfterItemUse() *event.Hook[*OnAfterItemUseEvent] {
	return u.onAfterItemUse
}

func (u *UserBase) OnNewLap() *event.Hook[*OnNewLapEvent] {
	return u.onNewLap
}

func (u *UserBase) OnBeforeNextStep() *event.Hook[*OnBeforeNextStepEvent] {
	return u.onBeforeNextStep
}

func (u *UserBase) OnAfterAction() *event.Hook[*OnAfterActionEvent] {
	return u.onAfterAction
}

func (u *UserBase) OnAfterMove() *event.Hook[*OnAfterMoveEvent] {
	return u.onAfterMove
}

func (u *UserBase) OnBeforeCurrentCell() *event.Hook[*OnBeforeCurrentCellEvent] {
	return u.onBeforeCurrentCell
}
