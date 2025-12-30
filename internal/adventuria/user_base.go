package adventuria

import (
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"sync"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type UserBase struct {
	core.BaseRecordProxy
	lastAction         *LastUserActionRecord
	inventory          Inventory
	timer              Timer
	stats              Stats
	inAction           bool
	mu                 sync.RWMutex
	hookIds            []string
	pEffectsUnsubGroup []event.UnsubGroup

	onAfterChooseGame   *event.Hook[*OnAfterChooseGameEvent]
	onAfterReroll       *event.Hook[*OnAfterRerollEvent]
	onBeforeDrop        *event.Hook[*OnBeforeDropEvent]
	onBeforeDropCheck   *event.Hook[*OnBeforeDropCheckEvent]
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
	onBeforeItemAdd     *event.Hook[*OnBeforeItemAdd]
	onAfterItemAdd      *event.Hook[*OnAfterItemAdd]
	onAfterItemSave     *event.Hook[*OnAfterItemSave]
}

func NewUser(userId string) (User, error) {
	var err error
	u := &UserBase{
		stats: Stats{},
	}

	err = u.fetchUser(userId)
	if err != nil {
		return nil, err
	}

	u.initHooks()

	u.timer, err = NewTimer(userId)
	if err != nil {
		return nil, err
	}

	u.lastAction, err = NewLastUserAction(u.Id)
	if err != nil {
		return nil, err
	}

	u.inventory, err = NewInventory(u, u.MaxInventorySlots())
	if err != nil {
		return nil, err
	}

	u.bindHooks()

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
	u.hookIds = make([]string, 3)

	u.hookIds[0] = PocketBase.OnRecordAfterUpdateSuccess(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			u.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	u.hookIds[1] = PocketBase.OnRecordUpdate(CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			e.Record.Set("stats", u.stats)
		}
		return e.Next()
	})

	i := 0
	u.pEffectsUnsubGroup = make([]event.UnsubGroup, len(persistentEffectsList))
	for _, effectCreator := range persistentEffectsList {
		effect := effectCreator()
		unsubs := effect.Subscribe(u)
		u.pEffectsUnsubGroup[i] = event.UnsubGroup{Fns: unsubs}
		i++
	}
	u.hookIds[2] = PocketBase.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		for _, unsubGroup := range u.pEffectsUnsubGroup {
			unsubGroup.Unsubscribe()
		}
		return e.Next()
	})
}

func (u *UserBase) Close() {
	PocketBase.OnRecordAfterUpdateSuccess(CollectionUsers).Unbind(u.hookIds[0])
	PocketBase.OnRecordUpdate(CollectionUsers).Unbind(u.hookIds[1])
	PocketBase.OnTerminate().Unbind(u.hookIds[2])
	for _, unsubGroup := range u.pEffectsUnsubGroup {
		unsubGroup.Unsubscribe()
	}
	u.inventory.Close()
	u.lastAction.Close()
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

func (u *UserBase) setCellsPassed(cellsPassed int) {
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

type CellReachedContext struct {
	User  User
	Moves []*OnAfterMoveEvent
}

func (u *UserBase) Move(steps int) ([]*OnAfterMoveEvent, error) {
	cellsPassed := u.CellsPassed()
	cellsCount := GameCells.Count()

	totalSteps := cellsPassed + steps
	currentCellNum := mod(totalSteps, cellsCount)
	lapsPassed := floorDiv(totalSteps, cellsCount) - floorDiv(cellsPassed, cellsCount)

	currentCell, ok := GameCells.GetByOrder(currentCellNum)
	if !ok {
		return nil, fmt.Errorf("user.Move(): cell with num = %d not found, steps = %d", currentCellNum, steps)
	}

	u.setCellsPassed(totalSteps)

	u.lastAction.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
	u.lastAction.SetUser(u.ID())
	u.lastAction.SetType(ActionTypeMove)
	u.lastAction.SetDiceRoll(steps)
	u.lastAction.setCell(currentCell.ID())
	u.lastAction.SetCanMove(false)

	onAfterMoveEvent := OnAfterMoveEvent{
		Steps:          steps,
		TotalSteps:     totalSteps,
		PrevTotalSteps: cellsPassed,
		CurrentCell:    currentCell,
		Laps:           lapsPassed,
	}

	res, err := u.OnAfterMove().Trigger(&onAfterMoveEvent)
	if res != nil && !res.Success {
		return nil, errors.New(res.Error)
	}
	if err != nil {
		return nil, err
	}

	cellReachedCtx := CellReachedContext{
		User:  u,
		Moves: []*OnAfterMoveEvent{&onAfterMoveEvent},
	}

	err = currentCell.OnCellReached(&cellReachedCtx)
	if err != nil {
		return nil, err
	}

	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		res, err = u.OnNewLap().Trigger(&OnNewLapEvent{
			Laps: lapsPassed,
		})
		if res != nil && !res.Success {
			return nil, errors.New(res.Error)
		}
		if err != nil {
			return nil, err
		}
	}

	return cellReachedCtx.Moves, nil
}

func (u *UserBase) CurrentCellOrder() int {
	return mod(u.CellsPassed(), GameCells.Count())
}

func (u *UserBase) MoveToClosestCellType(cellType CellType) ([]*OnAfterMoveEvent, error) {
	var (
		closest     int
		minDistance int
		found       bool
	)
	currentCellOrder := u.CurrentCellOrder()
	for order := range GameCells.GetOrderByType(cellType) {
		distance := abs(order - currentCellOrder)
		if !found || distance < minDistance {
			closest = order
			minDistance = distance
			found = true
		}
	}

	if !found {
		return nil, errors.New("cell not found")
	}

	return u.Move(closest - u.CurrentCellOrder())
}

func (u *UserBase) MoveToCellId(cellId string) ([]*OnAfterMoveEvent, error) {
	cellPos, ok := GameCells.GetOrderById(cellId)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellId)
	}
	return u.Move(cellPos - u.CurrentCellOrder())
}

func (u *UserBase) MoveToCellName(cellName string) ([]*OnAfterMoveEvent, error) {
	cellPos, ok := GameCells.GetOrderByName(cellName)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellName)
	}
	return u.Move(cellPos - u.CurrentCellOrder())
}

func (u *UserBase) MoveToClosestCellByNames(cellNames ...string) ([]*OnAfterMoveEvent, error) {
	if len(cellNames) == 0 {
		return nil, errors.New("moveToClosestCellByNames: cellNames is empty")
	}

	cellsOrder := make([]int, len(cellNames))
	for i, cellName := range cellNames {
		cellOrder, ok := GameCells.GetOrderByName(cellName)
		if !ok {
			return nil, errors.New("moveToClosestCellByNames: cell not found")
		}
		cellsOrder[i] = cellOrder
	}

	var (
		closest     int
		minDistance int
		found       bool
	)
	currentCellOrder := u.CurrentCellOrder()
	for order := range cellsOrder {
		distance := abs(order - currentCellOrder)
		if !found || distance < minDistance {
			closest = order
			minDistance = distance
			found = true
		}
	}

	if !found {
		return nil, errors.New("moveToClosestCellByNames: cell not found")
	}

	return u.Move(closest - u.CurrentCellOrder())
}

func (u *UserBase) Inventory() Inventory {
	return u.inventory
}

func (u *UserBase) LastAction() ActionRecord {
	return u.lastAction
}

func (u *UserBase) Timer() Timer {
	return u.timer
}

func (u *UserBase) Stats() *Stats {
	return &u.stats
}

func (u *UserBase) Balance() int {
	return u.GetInt("balance")
}

func (u *UserBase) SetBalance(balance int) {
	u.Set("balance", balance)
}

func (u *UserBase) IsStreamLive() bool {
	return u.GetBool("is_stream_live")
}

func (u *UserBase) SetIsStreamLive(isLive bool) {
	u.Set("is_stream_live", isLive)
}

func (u *UserBase) isInAction() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.inAction
}

func (u *UserBase) setIsInAction(b bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.inAction = b
}

func (u *UserBase) initHooks() {
	u.onAfterChooseGame = &event.Hook[*OnAfterChooseGameEvent]{}
	u.onAfterReroll = &event.Hook[*OnAfterRerollEvent]{}
	u.onBeforeDrop = &event.Hook[*OnBeforeDropEvent]{}
	u.onBeforeDropCheck = &event.Hook[*OnBeforeDropCheckEvent]{}
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
	u.onBeforeItemAdd = &event.Hook[*OnBeforeItemAdd]{}
	u.onAfterItemAdd = &event.Hook[*OnAfterItemAdd]{}
	u.onAfterItemSave = &event.Hook[*OnAfterItemSave]{}
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

func (u *UserBase) OnBeforeDropCheck() *event.Hook[*OnBeforeDropCheckEvent] {
	return u.onBeforeDropCheck
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

func (u *UserBase) OnBeforeItemAdd() *event.Hook[*OnBeforeItemAdd] {
	return u.onBeforeItemAdd
}

func (u *UserBase) OnAfterItemAdd() *event.Hook[*OnAfterItemAdd] {
	return u.onAfterItemAdd
}

func (u *UserBase) OnAfterItemSave() *event.Hook[*OnAfterItemSave] {
	return u.onAfterItemSave
}
