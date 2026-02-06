package adventuria

import (
	"adventuria/internal/adventuria/schema"
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

	balance int

	onAfterChooseGame   *event.Hook[*OnAfterChooseGameEvent]
	onAfterReroll       *event.Hook[*OnAfterRerollEvent]
	onBeforeDrop        *event.Hook[*OnBeforeDropEvent]
	onBeforeDropCheck   *event.Hook[*OnBeforeDropCheckEvent]
	onAfterDrop         *event.Hook[*OnAfterDropEvent]
	onAfterGoToJail     *event.Hook[*OnAfterGoToJailEvent]
	onBeforeDone        *event.Hook[*OnBeforeDoneEvent]
	onAfterDone         *event.Hook[*OnAfterDoneEvent]
	onBeforeRerollCheck *event.Hook[*OnBeforeRerollCheckEvent]
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
	onBeforeItemBuy     *event.Hook[*OnBeforeItemBuy]
	onBuyGetVariants    *event.Hook[*OnBuyGetVariants]
}

func NewUser(ctx AppContext, userId string) (User, error) {
	var err error
	u := &UserBase{
		stats: Stats{},
	}

	err = u.fetchUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	u.initHooks()

	u.timer, err = NewTimer(ctx, userId)
	if err != nil {
		return nil, err
	}

	u.lastAction, err = NewLastUserAction(ctx, u.Id)
	if err != nil {
		return nil, err
	}

	u.inventory, err = NewInventory(ctx, u, u.MaxInventorySlots())
	if err != nil {
		return nil, err
	}

	u.bindHooks(ctx)

	return u, nil
}

func NewUserFromName(ctx AppContext, name string) (User, error) {
	var record core.Record
	err := ctx.App.
		RecordQuery(schema.CollectionUsers).
		Where(dbx.HashExp{schema.UserSchema.Name: name}).
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return NewUser(ctx, record.Id)
}

func (u *UserBase) bindHooks(ctx AppContext) {
	u.hookIds = make([]string, 2)

	u.hookIds[0] = ctx.App.OnRecordUpdate(schema.CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			if e.Record.GetBool(schema.UserSchema.ClearStats) {
				e.Record.Set(schema.UserSchema.ClearStats, false)
				e.Record.Set(schema.UserSchema.Stats, "null")
			} else {
				e.Record.Set(schema.UserSchema.Stats, u.stats)
			}
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

	u.hookIds[1] = ctx.App.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		for _, unsubGroup := range u.pEffectsUnsubGroup {
			unsubGroup.Unsubscribe()
		}
		return e.Next()
	})
}

func (u *UserBase) Close(ctx AppContext) {
	ctx.App.OnRecordUpdate(schema.CollectionUsers).Unbind(u.hookIds[0])
	ctx.App.OnTerminate().Unbind(u.hookIds[1])
	for _, unsubGroup := range u.pEffectsUnsubGroup {
		unsubGroup.Unsubscribe()
	}
	u.inventory.Close(ctx)
	u.lastAction.Close(ctx)
}

func (u *UserBase) SetProxyRecord(record *core.Record) {
	u.BaseRecordProxy.SetProxyRecord(record)
	u.UnmarshalJSONField("stats", &u.stats)
	u.balance = u.GetInt(schema.UserSchema.Balance)
}

func (u *UserBase) fetchUser(ctx AppContext, userId string) error {
	user, err := ctx.App.FindRecordById(schema.CollectionUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)

	return nil
}

func (u *UserBase) Refetch(ctx AppContext) error {
	return u.fetchUser(ctx, u.Id)
}

func (u *UserBase) ID() string {
	return u.Id
}

func (u *UserBase) Name() string {
	return u.GetString(schema.UserSchema.Name)
}

func (u *UserBase) IsSafeDrop() bool {
	return u.DropsInARow() < GameSettings.DropsToJail()
}

func (u *UserBase) IsInJail() bool {
	return u.GetBool(schema.UserSchema.IsInJail)
}

func (u *UserBase) SetIsInJail(b bool) {
	u.Set(schema.UserSchema.IsInJail, b)
}

func (u *UserBase) CurrentCell() (Cell, bool) {
	currentCellNum := u.CellsPassed() % GameCells.Count()
	cell, ok := GameCells.GetByOrder(currentCellNum)

	return cell, ok
}

func (u *UserBase) Points() int {
	return u.GetInt(schema.UserSchema.Points)
}

func (u *UserBase) SetPoints(points int) {
	u.Set(schema.UserSchema.Points, points)
}

func (u *UserBase) DropsInARow() int {
	return u.GetInt(schema.UserSchema.DropsInARow)
}

func (u *UserBase) SetDropsInARow(drops int) {
	u.Set(schema.UserSchema.DropsInARow, drops)
}

func (u *UserBase) CellsPassed() int {
	return u.GetInt(schema.UserSchema.CellsPassed)
}

func (u *UserBase) setCellsPassed(cellsPassed int) {
	u.Set(schema.UserSchema.CellsPassed, cellsPassed)
}

func (u *UserBase) MaxInventorySlots() int {
	return u.GetInt(schema.UserSchema.MaxInventorySlots)
}

func (u *UserBase) SetMaxInventorySlots(maxInventorySlots int) {
	u.Set(schema.UserSchema.MaxInventorySlots, maxInventorySlots)
}

func (u *UserBase) ItemWheelsCount() int {
	return u.GetInt(schema.UserSchema.ItemWheelsCount)
}

func (u *UserBase) SetItemWheelsCount(itemWheelsCount int) {
	u.Set(schema.UserSchema.ItemWheelsCount, itemWheelsCount)
}

func (u *UserBase) Move(ctx AppContext, steps int) ([]*MoveResult, error) {
	prevCell, ok := u.CurrentCell()
	if ok {
		err := prevCell.OnCellLeft(&CellLeftContext{
			AppContext: ctx,
			User:       u,
		})
		if err != nil {
			return nil, err
		}
	}

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
	if err := ctx.App.Save(u.ProxyRecord()); err != nil {
		return nil, err
	}

	u.lastAction.SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionActions)))
	u.lastAction.SetUser(u.ID())
	u.lastAction.SetType(ActionTypeMove)
	u.lastAction.SetDiceRoll(steps)
	u.lastAction.setCell(currentCell.ID())
	u.lastAction.SetCanMove(false)

	if err := ctx.App.Save(u.lastAction.ProxyRecord()); err != nil {
		return nil, err
	}

	onAfterMoveEvent := OnAfterMoveEvent{
		AppContext:     ctx,
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
		AppContext: ctx,
		User:       u,
		Moves: []*MoveResult{
			{
				Steps:          onAfterMoveEvent.Steps,
				TotalSteps:     onAfterMoveEvent.TotalSteps,
				PrevTotalSteps: onAfterMoveEvent.PrevTotalSteps,
				CurrentCell:    onAfterMoveEvent.CurrentCell,
				Laps:           onAfterMoveEvent.Laps,
			},
		},
	}

	err = currentCell.OnCellReached(&cellReachedCtx)
	if err != nil {
		return nil, err
	}

	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		res, err = u.OnNewLap().Trigger(&OnNewLapEvent{
			AppContext: ctx,
			Laps:       lapsPassed,
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

func (u *UserBase) MoveToClosestCellType(ctx AppContext, cellType CellType) ([]*MoveResult, error) {
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

	return u.Move(ctx, closest-u.CurrentCellOrder())
}

func (u *UserBase) MoveToCellId(ctx AppContext, cellId string) ([]*MoveResult, error) {
	cellPos, ok := GameCells.GetOrderById(cellId)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellId)
	}
	return u.Move(ctx, cellPos-u.CurrentCellOrder())
}

func (u *UserBase) MoveToCellName(ctx AppContext, cellName string) ([]*MoveResult, error) {
	cellPos, ok := GameCells.GetOrderByName(cellName)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellName)
	}
	return u.Move(ctx, cellPos-u.CurrentCellOrder())
}

func (u *UserBase) MoveToClosestCellByNames(ctx AppContext, cellNames ...string) ([]*MoveResult, error) {
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
	for _, order := range cellsOrder {
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

	return u.Move(ctx, closest-u.CurrentCellOrder())
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
	return u.balance
}

func (u *UserBase) AddBalance(ctx AppContext, amount int) error {
	query := fmt.Sprintf(
		"UPDATE %s SET %[2]s = %[2]s + {:amount} WHERE id = {:id}",
		schema.CollectionUsers,
		schema.UserSchema.Balance,
	)
	_, err := ctx.App.DB().NewQuery(query).Bind(dbx.Params{
		"amount": amount,
		"id":     u.ID(),
	}).Execute()
	if err != nil {
		return err
	}

	u.balance += amount

	return nil
}

func (u *UserBase) IsStreamLive() bool {
	return u.GetBool(schema.UserSchema.IsStreamLive)
}

func (u *UserBase) SetIsStreamLive(isLive bool) {
	u.Set(schema.UserSchema.IsStreamLive, isLive)
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
	u.onBeforeRerollCheck = &event.Hook[*OnBeforeRerollCheckEvent]{}
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
	u.onBeforeItemBuy = &event.Hook[*OnBeforeItemBuy]{}
	u.onBuyGetVariants = &event.Hook[*OnBuyGetVariants]{}
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

func (u *UserBase) OnBeforeRerollCheck() *event.Hook[*OnBeforeRerollCheckEvent] {
	return u.onBeforeRerollCheck
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

func (u *UserBase) OnBeforeItemBuy() *event.Hook[*OnBeforeItemBuy] {
	return u.onBeforeItemBuy
}

func (u *UserBase) OnBuyGetVariants() *event.Hook[*OnBuyGetVariants] {
	return u.onBuyGetVariants
}
