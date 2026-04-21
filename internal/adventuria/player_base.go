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

type PlayerBase struct {
	core.BaseRecordProxy
	lastAction         *LastPlayerActionRecord
	inventory          Inventory
	progress           PlayerProgress
	locked             bool
	mu                 sync.RWMutex
	hookIds            []string
	pEffectsUnsubGroup []event.UnsubGroup

	onAfterChooseGame      *event.Hook[*OnAfterChooseGameEvent]
	onAfterReroll          *event.Hook[*OnAfterRerollEvent]
	onBeforeDrop           *event.Hook[*OnBeforeDropEvent]
	onBeforeDropCheck      *event.Hook[*OnBeforeDropCheckEvent]
	onAfterDrop            *event.Hook[*OnAfterDropEvent]
	onAfterGoToJail        *event.Hook[*OnAfterGoToJailEvent]
	onBeforeDone           *event.Hook[*OnBeforeDoneEvent]
	onAfterDone            *event.Hook[*OnAfterDoneEvent]
	onBeforeRerollCheck    *event.Hook[*OnBeforeRerollCheckEvent]
	onBeforeRoll           *event.Hook[*OnBeforeRollEvent]
	onBeforeRollMove       *event.Hook[*OnBeforeRollMoveEvent]
	onAfterRoll            *event.Hook[*OnAfterRollEvent]
	onBeforeWheelRoll      *event.Hook[*OnBeforeWheelRollEvent]
	onAfterWheelRoll       *event.Hook[*OnAfterWheelRollEvent]
	onAfterItemRoll        *event.Hook[*OnAfterItemRollEvent]
	onAfterItemUse         *event.Hook[*OnAfterItemUseEvent]
	onNewLap               *event.Hook[*OnNewLapEvent]
	onBeforeNextStep       *event.Hook[*OnBeforeNextStepEvent]
	onAfterAction          *event.Hook[*OnAfterActionEvent]
	onAfterMove            *event.Hook[*OnAfterMoveEvent]
	onBeforeCurrentCell    *event.Hook[*OnBeforeCurrentCellEvent]
	onBeforeItemAdd        *event.Hook[*OnBeforeItemAdd]
	onAfterItemAdd         *event.Hook[*OnAfterItemAdd]
	onAfterItemSave        *event.Hook[*OnAfterItemSave]
	onBeforeItemBuy        *event.Hook[*OnBeforeItemBuy]
	onBuyGetVariants       *event.Hook[*OnBuyGetVariants]
	onBeforeTeleportOnCell *event.Hook[*OnBeforeTeleportOnCell]
}

func NewPlayer(ctx AppContext, playerId string) (Player, error) {
	var err error
	u := &PlayerBase{}

	err = u.fetchPlayer(ctx, playerId)
	if err != nil {
		return nil, err
	}

	u.initHooks()

	u.progress, err = NewPlayerProgress(ctx, playerId, GameSettings.CurrentSeason())
	if err != nil {
		return nil, err
	}

	u.lastAction, err = NewLastPlayerAction(ctx, u.Id)
	if err != nil {
		return nil, err
	}

	u.inventory, err = NewInventory(ctx, u)
	if err != nil {
		return nil, err
	}

	u.bindHooks(ctx)

	return u, nil
}

func NewPlayerFromName(ctx AppContext, playerName string) (Player, error) {
	var record core.Record
	err := ctx.App.
		RecordQuery(schema.CollectionPlayers).
		Where(dbx.HashExp{schema.PlayerSchema.Name: playerName}).
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return NewPlayer(ctx, record.Id)
}

func (p *PlayerBase) bindHooks(ctx AppContext) {
	p.hookIds = make([]string, 1)

	i := 0
	p.pEffectsUnsubGroup = make([]event.UnsubGroup, len(persistentEffectsList))
	for _, effectCreator := range persistentEffectsList {
		effect := effectCreator()
		unsubs := effect.Subscribe(p)
		p.pEffectsUnsubGroup[i] = event.UnsubGroup{Fns: unsubs}
		i++
	}

	p.hookIds[0] = ctx.App.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		for _, unsubGroup := range p.pEffectsUnsubGroup {
			unsubGroup.Unsubscribe()
		}
		return e.Next()
	})
}

func (p *PlayerBase) Close(ctx AppContext) {
	ctx.App.OnTerminate().Unbind(p.hookIds[0])
	for _, unsubGroup := range p.pEffectsUnsubGroup {
		unsubGroup.Unsubscribe()
	}
	p.progress.Close(ctx)
	p.inventory.Close(ctx)
}

func (p *PlayerBase) fetchPlayer(ctx AppContext, playerId string) error {
	player, err := ctx.App.FindRecordById(schema.CollectionPlayers, playerId)
	if err != nil {
		return err
	}

	p.SetProxyRecord(player)

	return nil
}

func (p *PlayerBase) Refetch(ctx AppContext) error {
	if err := p.fetchPlayer(ctx, p.Id); err != nil {
		return err
	}
	if err := p.progress.Refetch(ctx); err != nil {
		return err
	}
	if err := p.lastAction.Refetch(ctx); err != nil {
		return err
	}
	return p.inventory.Refetch(ctx)
}

func (p *PlayerBase) ID() string {
	return p.Id
}

func (p *PlayerBase) Name() string {
	return p.GetString(schema.PlayerSchema.Name)
}

func (p *PlayerBase) Move(ctx AppContext, steps int) ([]*MoveResult, error) {
	prevCell, ok := p.progress.CurrentCell()
	if ok {
		err := prevCell.OnCellLeft(&CellLeftContext{
			AppContext: ctx,
			Player:     p,
		})
		if err != nil {
			return nil, err
		}
	}

	cellsPassed := p.progress.CellsPassed()
	cellsCount := GameCells.Count()

	totalSteps := cellsPassed + steps
	currentCellNum := mod(totalSteps, cellsCount)
	lapsPassed := floorDiv(totalSteps, cellsCount) - floorDiv(cellsPassed, cellsCount)

	currentCell, ok := GameCells.GetByOrder(currentCellNum)
	if !ok {
		return nil, fmt.Errorf("player.Move(): cell with num = %d not found, steps = %d", currentCellNum, steps)
	}

	p.progress.addCellsPassed(steps)

	p.lastAction.SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionActions)))
	p.lastAction.SetPlayer(p.ID())
	p.lastAction.SetType(ActionTypeMove)
	p.lastAction.SetCellsPassed(steps)
	p.lastAction.setCell(currentCell.ID())
	p.lastAction.SetCanMove(false)

	if err := ctx.App.Save(p.lastAction.ProxyRecord()); err != nil {
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

	res, err := p.OnAfterMove().Trigger(&onAfterMoveEvent)
	if res != nil && !res.Success {
		return nil, errors.New(res.Error)
	}
	if err != nil {
		return nil, err
	}

	cellReachedCtx := CellReachedContext{
		AppContext: ctx,
		Player:     p,
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
		res, err = p.OnNewLap().Trigger(&OnNewLapEvent{
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

func (p *PlayerBase) MoveToClosestCellType(ctx AppContext, cellType CellType) ([]*MoveResult, error) {
	var (
		closest     int
		minDistance int
		found       bool
	)
	currentCellOrder := p.progress.CurrentCellOrder()
	for order := range GameCells.GetOrderByType(cellType) {
		distance := abs(order - currentCellOrder)
		if !found ||
			distance < minDistance ||
			(distance == minDistance && order > closest) {
			closest = order
			minDistance = distance
			found = true
		}
	}

	if !found {
		return nil, errors.New("cell not found")
	}

	return p.Move(ctx, closest-p.progress.CurrentCellOrder())
}

func (p *PlayerBase) MoveToCellId(ctx AppContext, cellId string) ([]*MoveResult, error) {
	cellPos, ok := GameCells.GetOrderById(cellId)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellId)
	}
	return p.Move(ctx, cellPos-p.progress.CurrentCellOrder())
}

func (p *PlayerBase) MoveToCellName(ctx AppContext, cellName string) ([]*MoveResult, error) {
	cellPos, ok := GameCells.GetOrderByName(cellName)
	if !ok {
		return nil, fmt.Errorf("cell %s not found", cellName)
	}
	return p.Move(ctx, cellPos-p.progress.CurrentCellOrder())
}

func (p *PlayerBase) MoveToCellOrder(ctx AppContext, cellOrder int) ([]*MoveResult, error) {
	return p.Move(ctx, cellOrder-p.progress.CurrentCellOrder())
}

func (p *PlayerBase) MoveToClosestCellByNames(ctx AppContext, cellNames ...string) ([]*MoveResult, error) {
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
	currentCellOrder := p.progress.CurrentCellOrder()
	for _, order := range cellsOrder {
		distance := abs(order - currentCellOrder)
		if !found ||
			distance < minDistance ||
			(distance == minDistance && order > closest) {
			closest = order
			minDistance = distance
			found = true
		}
	}

	if !found {
		return nil, errors.New("moveToClosestCellByNames: cell not found")
	}

	return p.Move(ctx, closest-p.progress.CurrentCellOrder())
}

func (p *PlayerBase) Inventory() Inventory {
	return p.inventory
}

func (p *PlayerBase) LastAction() ActionRecord {
	return p.lastAction
}

func (p *PlayerBase) Progress() PlayerProgress {
	return p.progress
}

func (p *PlayerBase) IsStreamLive() bool {
	return p.GetBool(schema.PlayerSchema.IsStreamLive)
}

func (p *PlayerBase) SetIsStreamLive(isLive bool) {
	p.Set(schema.PlayerSchema.IsStreamLive, isLive)
}

func (p *PlayerBase) Locked() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.locked
}

func (p *PlayerBase) Lock() {
	p.mu.Lock()
	p.locked = true
}

func (p *PlayerBase) Unlock() {
	p.mu.Unlock()
	p.locked = false
}

func (p *PlayerBase) initHooks() {
	p.onAfterChooseGame = &event.Hook[*OnAfterChooseGameEvent]{}
	p.onAfterReroll = &event.Hook[*OnAfterRerollEvent]{}
	p.onBeforeDrop = &event.Hook[*OnBeforeDropEvent]{}
	p.onBeforeDropCheck = &event.Hook[*OnBeforeDropCheckEvent]{}
	p.onAfterDrop = &event.Hook[*OnAfterDropEvent]{}
	p.onAfterGoToJail = &event.Hook[*OnAfterGoToJailEvent]{}
	p.onBeforeDone = &event.Hook[*OnBeforeDoneEvent]{}
	p.onAfterDone = &event.Hook[*OnAfterDoneEvent]{}
	p.onBeforeRerollCheck = &event.Hook[*OnBeforeRerollCheckEvent]{}
	p.onBeforeRoll = &event.Hook[*OnBeforeRollEvent]{}
	p.onBeforeRollMove = &event.Hook[*OnBeforeRollMoveEvent]{}
	p.onAfterRoll = &event.Hook[*OnAfterRollEvent]{}
	p.onBeforeWheelRoll = &event.Hook[*OnBeforeWheelRollEvent]{}
	p.onAfterWheelRoll = &event.Hook[*OnAfterWheelRollEvent]{}
	p.onAfterItemRoll = &event.Hook[*OnAfterItemRollEvent]{}
	p.onAfterItemUse = &event.Hook[*OnAfterItemUseEvent]{}
	p.onNewLap = &event.Hook[*OnNewLapEvent]{}
	p.onBeforeNextStep = &event.Hook[*OnBeforeNextStepEvent]{}
	p.onAfterAction = &event.Hook[*OnAfterActionEvent]{}
	p.onAfterMove = &event.Hook[*OnAfterMoveEvent]{}
	p.onBeforeCurrentCell = &event.Hook[*OnBeforeCurrentCellEvent]{}
	p.onBeforeItemAdd = &event.Hook[*OnBeforeItemAdd]{}
	p.onAfterItemAdd = &event.Hook[*OnAfterItemAdd]{}
	p.onAfterItemSave = &event.Hook[*OnAfterItemSave]{}
	p.onBeforeItemBuy = &event.Hook[*OnBeforeItemBuy]{}
	p.onBuyGetVariants = &event.Hook[*OnBuyGetVariants]{}
	p.onBeforeTeleportOnCell = &event.Hook[*OnBeforeTeleportOnCell]{}
}

func (p *PlayerBase) OnAfterChooseGame() *event.Hook[*OnAfterChooseGameEvent] {
	return p.onAfterChooseGame
}

func (p *PlayerBase) OnAfterReroll() *event.Hook[*OnAfterRerollEvent] {
	return p.onAfterReroll
}

func (p *PlayerBase) OnBeforeDrop() *event.Hook[*OnBeforeDropEvent] {
	return p.onBeforeDrop
}

func (p *PlayerBase) OnBeforeDropCheck() *event.Hook[*OnBeforeDropCheckEvent] {
	return p.onBeforeDropCheck
}

func (p *PlayerBase) OnAfterDrop() *event.Hook[*OnAfterDropEvent] {
	return p.onAfterDrop
}

func (p *PlayerBase) OnAfterGoToJail() *event.Hook[*OnAfterGoToJailEvent] {
	return p.onAfterGoToJail
}

func (p *PlayerBase) OnBeforeDone() *event.Hook[*OnBeforeDoneEvent] {
	return p.onBeforeDone
}

func (p *PlayerBase) OnAfterDone() *event.Hook[*OnAfterDoneEvent] {
	return p.onAfterDone
}

func (p *PlayerBase) OnBeforeRerollCheck() *event.Hook[*OnBeforeRerollCheckEvent] {
	return p.onBeforeRerollCheck
}

func (p *PlayerBase) OnBeforeRoll() *event.Hook[*OnBeforeRollEvent] {
	return p.onBeforeRoll
}

func (p *PlayerBase) OnBeforeRollMove() *event.Hook[*OnBeforeRollMoveEvent] {
	return p.onBeforeRollMove
}

func (p *PlayerBase) OnAfterRoll() *event.Hook[*OnAfterRollEvent] {
	return p.onAfterRoll
}

func (p *PlayerBase) OnBeforeWheelRoll() *event.Hook[*OnBeforeWheelRollEvent] {
	return p.onBeforeWheelRoll
}

func (p *PlayerBase) OnAfterWheelRoll() *event.Hook[*OnAfterWheelRollEvent] {
	return p.onAfterWheelRoll
}

func (p *PlayerBase) OnAfterItemRoll() *event.Hook[*OnAfterItemRollEvent] {
	return p.onAfterItemRoll
}

func (p *PlayerBase) OnAfterItemUse() *event.Hook[*OnAfterItemUseEvent] {
	return p.onAfterItemUse
}

func (p *PlayerBase) OnNewLap() *event.Hook[*OnNewLapEvent] {
	return p.onNewLap
}

func (p *PlayerBase) OnBeforeNextStep() *event.Hook[*OnBeforeNextStepEvent] {
	return p.onBeforeNextStep
}

func (p *PlayerBase) OnAfterAction() *event.Hook[*OnAfterActionEvent] {
	return p.onAfterAction
}

func (p *PlayerBase) OnAfterMove() *event.Hook[*OnAfterMoveEvent] {
	return p.onAfterMove
}

func (p *PlayerBase) OnBeforeCurrentCell() *event.Hook[*OnBeforeCurrentCellEvent] {
	return p.onBeforeCurrentCell
}

func (p *PlayerBase) OnBeforeItemAdd() *event.Hook[*OnBeforeItemAdd] {
	return p.onBeforeItemAdd
}

func (p *PlayerBase) OnAfterItemAdd() *event.Hook[*OnAfterItemAdd] {
	return p.onAfterItemAdd
}

func (p *PlayerBase) OnAfterItemSave() *event.Hook[*OnAfterItemSave] {
	return p.onAfterItemSave
}

func (p *PlayerBase) OnBeforeItemBuy() *event.Hook[*OnBeforeItemBuy] {
	return p.onBeforeItemBuy
}

func (p *PlayerBase) OnBuyGetVariants() *event.Hook[*OnBuyGetVariants] {
	return p.onBuyGetVariants
}

func (p *PlayerBase) OnBeforeTeleportOnCell() *event.Hook[*OnBeforeTeleportOnCell] {
	return p.onBeforeTeleportOnCell
}
