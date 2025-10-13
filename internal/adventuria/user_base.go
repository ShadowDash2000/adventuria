package adventuria

import (
	"adventuria/pkg/event"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type UserBase struct {
	core.BaseRecordProxy
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

func NewUser(userId string) (User, error) {
	if userId == "" {
		return nil, errors.New("empty user id")
	}

	var err error
	timer, err := NewTimer(userId)
	if err != nil {
		return nil, err
	}

	u := &UserBase{
		timer: timer,
	}

	err = u.fetchUser(userId)
	if err != nil {
		return nil, err
	}

	u.lastAction, err = NewLastUserAction(userId)
	if err != nil {
		return nil, err
	}

	u.inventory, err = NewInventory(u, u.MaxInventorySlots())
	if err != nil {
		return nil, err
	}

	u.bindHooks()
	u.initHooks()

	return u, nil
}

func NewUserFromName(name string) (User, error) {
	record, err := PocketBase.FindRecordsByFilter(
		CollectionUsers,
		"name = {:name}",
		"",
		1,
		0,
		dbx.Params{
			"name": name,
		},
	)
	if err != nil {
		return nil, err
	}

	return NewUser(record[0].Id)
}

func (u *UserBase) bindHooks() {
	PocketBase.OnRecordAfterUpdateSuccess(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
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
	u.UnmarshalJSONField("stats", &u.stats)
}

func (u *UserBase) fetchUser(userId string) error {
	user, err := PocketBase.FindRecordById(CollectionUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)

	return nil
}

func (u *UserBase) ID() string {
	return u.Id
}

func (u *UserBase) Name() string {
	return u.GetString("name")
}

func (u *UserBase) IsSafeDrop() bool {
	return u.DropsInARow() < GameSettings.DropsToJail()
}

func (u *UserBase) IsInJail() bool {
	return u.GetBool("isInJail")
}

func (u *UserBase) SetIsInJail(b bool) {
	u.Set("isInJail", b)
}

func (u *UserBase) CurrentCell() (Cell, bool) {
	currentCellNum := u.CellsPassed() % GameCells.Count()
	cell, ok := GameCells.GetByOrder(currentCellNum)

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

	return PocketBase.Save(u)
}

func (u *UserBase) Move(steps int) (*OnAfterMoveEvent, error) {
	cellsPassed := u.CellsPassed()
	currentCellNum := (cellsPassed + steps) % GameCells.Count()
	currentCell, _ := GameCells.GetByOrder(currentCellNum)

	u.SetCellsPassed(cellsPassed + steps)

	err := currentCell.OnCellReached(u)
	if err != nil {
		return nil, err
	}

	prevCellNum := cellsPassed % GameCells.Count()
	lapsPassed := (prevCellNum + steps) / GameCells.Count()
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
	cellPos, ok := GameCells.GetOrderByType(cellType)
	if !ok {
		return errors.New("cell not found")
	}

	currentCellNum := u.CellsPassed() % GameCells.Count()

	_, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserBase) MoveToCellId(cellId string) error {
	cellPos, ok := GameCells.GetOrderById(cellId)
	if !ok {
		return fmt.Errorf("cell %s not found", cellId)
	}

	currentCellNum := u.CellsPassed() % GameCells.Count()

	_, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

// NextAction
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *UserBase) NextAction() ActionType {
	if u.LastAction() == nil {
		// TODO change hardcoded value
		return "rollDice"
	}

	return u.LastAction().NextAction()
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

func (u *UserBase) initHooks() {
	u.onAfterChooseGame = &event.Hook[*OnAfterChooseGameEvent]{}
	u.onAfterReroll = &event.Hook[*OnAfterRerollEvent]{}
	u.onBeforeDrop = &event.Hook[*OnBeforeDropEvent]{}
	u.onAfterDrop = &event.Hook[*OnAfterDropEvent]{}
	u.onAfterGoToJail = &event.Hook[*OnAfterGoToJailEvent]{}
	u.onBeforeDone = &event.Hook[*OnBeforeDoneEvent]{}
	u.onAfterDone = &event.Hook[*OnAfterDoneEvent]{}
	u.onBeforeRoll = &event.Hook[*OnBeforeRollEvent]{}
	u.onBeforeRollMove = &event.Hook[*OnBeforeRollMoveEvent]{}
	u.onAfterRoll = &event.Hook[*OnAfterRollEvent]{}
	u.onBeforeWheelRoll = &event.Hook[*OnBeforeWheelRollEvent]{}
	u.onAfterWheelRoll = &event.Hook[*OnAfterWheelRollEvent]{}
	u.onAfterItemRoll = &event.Hook[*OnAfterItemRollEvent]{}
	u.onAfterItemUse = &event.Hook[*OnAfterItemUseEvent]{}
	u.onNewLap = &event.Hook[*OnNewLapEvent]{}
	u.onBeforeNextStep = &event.Hook[*OnBeforeNextStepEvent]{}
	u.onAfterAction = &event.Hook[*OnAfterActionEvent]{}
	u.onAfterMove = &event.Hook[*OnAfterMoveEvent]{}
	u.onBeforeCurrentCell = &event.Hook[*OnBeforeCurrentCellEvent]{}
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
