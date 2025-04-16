package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"adventuria/pkg/helper"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

var (
	GameApp         core.App
	GameCells       *Cells
	GameItems       *Items
	GameCollections *collections.Collections
	GameSettings    *Settings
	GameEvent       Event
)

type BaseGame struct {
	app   core.App
	users *cache.MemoryCache[string, *User]
}

func New(app core.App) Game {
	return &BaseGame{
		app:   app,
		users: cache.NewMemoryCache[string, *User](0, true),
	}
}

func (g *BaseGame) Init() {
	GameApp = g.app
	GameCells = NewCells()
	GameItems = NewItems()
	GameCollections = collections.NewCollections(g.app)
	GameSettings = NewSettings()
	GameEvent = NewEvent()
}

func (g *BaseGame) GetUser(userId string) (*User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *BaseGame) afterAction(user *User, event EffectUse) error {
	onAfterActionFields := &OnAfterActionFields{
		Event: event,
	}
	GameEvent.Go(OnAfterAction, NewEventFields(user, onAfterActionFields))

	err := user.Save()
	if err != nil {
		return err
	}

	_, err = user.Inventory.ApplyEffectsByEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) GetNextStepType(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	return nextStepType, nil
}

func (g *BaseGame) UpdateAction(actionId string, comment string, file *filesystem.File, userId string) error {
	actionsCollection, err := GameCollections.Get(TableActions)
	if err != nil {
		return err
	}

	record := &core.Record{}
	err = GameApp.
		RecordQuery(actionsCollection).
		AndWhere(
			dbx.HashExp{
				"user": userId,
				"id":   actionId,
			},
		).
		AndWhere(
			dbx.Or(
				dbx.HashExp{"type": ActionTypeDone},
				dbx.HashExp{"type": ActionTypeDrop},
				dbx.HashExp{"type": ActionTypeReroll},
			),
		).
		Limit(1).
		One(record)
	if err != nil {
		return err
	}

	action := NewActionFromRecord(record)
	action.SetComment(comment)

	return action.Save()
}

func (g *BaseGame) Reroll(comment string, file *filesystem.File, userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != ActionTypeChooseResult {
		return errors.New("next step isn't choose result")
	}

	currentCell, _ := user.CurrentCell()

	if currentCell.CantReroll() {
		return errors.New("can't reroll on this cell")
	}

	action := user.LastAction
	action.SetType(ActionTypeReroll)
	action.SetComment(comment)
	err = action.Save()
	if err != nil {
		return err
	}

	onAfterRerollFields := &OnAfterRerollFields{}
	GameEvent.Go(OnAfterReroll, NewEventFields(user, onAfterRerollFields))

	err = g.afterAction(user, EffectUseOnReroll)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) RollDice(userId string) (int, []int, Cell, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, nil, nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return 0, nil, nil, err
	}

	if nextStepType != ActionTypeRollDice {
		return 0, nil, nil, errors.New("next step isn't roll dice")
	}

	onBeforeRollFields := &OnBeforeRollFields{
		Dices: []Dice{DiceTypeD6, DiceTypeD6},
	}
	GameEvent.Go(OnBeforeRoll, NewEventFields(user, onBeforeRollFields))

	onBeforeRollMoveFields := &OnBeforeRollMoveFields{
		N: 0,
	}
	diceRolls := make([]int, len(onBeforeRollFields.Dices))
	for i, dice := range onBeforeRollFields.Dices {
		diceRolls[i] = dice.Roll()
		onBeforeRollMoveFields.N += diceRolls[i]
	}

	GameEvent.Go(OnBeforeRollMove, NewEventFields(user, onBeforeRollMoveFields))

	fields, err := user.Move(onBeforeRollMoveFields.N)

	onAfterRollFields := &OnAfterRollFields{
		Dices: onBeforeRollFields.Dices,
		N:     onBeforeRollMoveFields.N,
	}
	GameEvent.Go(OnAfterRoll, NewEventFields(user, onAfterRollFields))

	err = g.afterAction(user, EffectUseOnRoll)
	if err != nil {
		return 0, nil, nil, err
	}

	return onBeforeRollMoveFields.N, diceRolls, fields.CurrentCell, nil
}

func (g *BaseGame) Drop(comment string, file *filesystem.File, userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != ActionTypeChooseResult {
		return errors.New("next step isn't choose result")
	}

	currentCell, _ := user.CurrentCell()

	if currentCell.CantDrop() || !user.CanDrop() {
		return errors.New("can't drop on this cell")
	}

	onBeforeDropFields := &OnBeforeDropFields{
		IsSafeDrop: false,
	}
	GameEvent.Go(OnBeforeDrop, NewEventFields(user, onBeforeDropFields))

	action := user.LastAction
	action.SetType(ActionTypeDrop)
	action.SetComment(comment)
	err = action.Save()
	if err != nil {
		return err
	}

	if !onBeforeDropFields.IsSafeDrop && !currentCell.IsSafeDrop() {
		user.SetPoints(user.Points() + GameSettings.PointsForDrop())
		user.SetDropsInARow(user.DropsInARow() + 1)

		if !user.IsSafeDrop() {
			if err = user.MoveToJail(); err != nil {
				return err
			}
		}
	}

	onAfterDropFields := &OnAfterDropFields{}
	GameEvent.Go(OnAfterDrop, NewEventFields(user, onAfterDropFields))

	err = g.afterAction(user, EffectUseOnDrop)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) Done(comment string, file *filesystem.File, userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != ActionTypeChooseResult {
		return errors.New("next step isn't choose result")
	}

	onBeforeDoneFields := &OnBeforeDoneFields{
		CellPointsDivide: 0,
	}
	GameEvent.Go(OnBeforeDone, NewEventFields(user, onBeforeDoneFields))

	currentCell, _ := user.CurrentCell()

	action := user.LastAction
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	err = action.Save()
	if err != nil {
		return err
	}

	cellPoints := currentCell.Points()
	if onBeforeDoneFields.CellPointsDivide != 0 {
		cellPoints /= onBeforeDoneFields.CellPointsDivide
	}

	user.SetDropsInARow(0)
	user.SetIsInJail(false)
	user.SetCantDrop(false)
	user.SetPoints(user.Points() + cellPoints)

	onAfterDoneFields := &OnAfterDoneFields{}
	GameEvent.Go(OnAfterDone, NewEventFields(user, onAfterDoneFields))

	err = g.afterAction(user, EffectUseOnDone)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) RollItem(userId string) (*WheelRollResult, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	if user.ItemWheelsCount() <= 0 {
		return nil, errors.New("no item wheels left")
	}

	itemsCol, err := GameCollections.Get(TableItems)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: itemsCol,
		EffectUse:  EffectUseOnRollItem,
	}

	items := GameItems.GetAllRollable()

	if len(items) == 0 {
		return nil, errors.New("items not found")
	}

	for _, item := range items {
		res.FillerItems = append(res.FillerItems, &WheelItem{
			Name: item.Name(),
			Icon: item.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()

	err = user.Inventory.MustAddItemById(res.WinnerId)
	if err != nil {
		return nil, err
	}

	user.SetItemWheelsCount(user.ItemWheelsCount() - 1)

	onAfterItemRollFields := &OnAfterItemRollFields{
		ItemId: res.WinnerId,
	}
	GameEvent.Go(OnAfterItemRoll, NewEventFields(user, onAfterItemRollFields))
	onAfterWheelRollFields := &OnAfterWheelRollFields{
		ItemId: res.WinnerId,
	}
	GameEvent.Go(OnAfterWheelRoll, NewEventFields(user, onAfterWheelRollFields))

	err = g.afterAction(user, EffectUseOnRollItem)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *BaseGame) RollWheel(userId string) (*WheelRollResult, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return nil, err
	}

	if nextStepType != ActionTypeRollWheel {
		return nil, errors.New("next step isn't roll wheel")
	}

	currentCell, _ := user.CurrentCell()
	onBeforeWheelRollFields := &OnBeforeWheelRollFields{
		CurrentCell: currentCell.(CellWheel),
	}
	GameEvent.Go(OnBeforeWheelRoll, NewEventFields(user, onBeforeWheelRollFields))

	res, err := onBeforeWheelRollFields.CurrentCell.Roll(user)
	if err != nil {
		return nil, err
	}

	action := user.LastAction
	action.SetType(ActionTypeRollWheel)
	action.SetValue(res.WinnerId)
	action.SetCollectionRef(res.Collection.Id)
	err = action.Save()
	if err != nil {
		return nil, err
	}

	onAfterWheelRollFields := &OnAfterWheelRollFields{
		ItemId: res.WinnerId,
	}
	GameEvent.Go(OnAfterWheelRoll, NewEventFields(user, onAfterWheelRollFields))

	err = g.afterAction(user, res.EffectUse)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *BaseGame) GetLastAction(userId string) (bool, Action, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return false, nil, err
	}

	return user.CanDrop(), user.LastAction, nil
}

func (g *BaseGame) GetItemsEffects(userId string, event EffectUse) (*Effects, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	effects, _, err := user.Inventory.Effects(event)
	if err != nil {
		return nil, err
	}

	return effects, nil
}

func (g *BaseGame) UseItem(userId, itemId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	err = user.Inventory.UseItem(itemId)
	if err != nil {
		return err
	}

	onAfterItemUseFields := &OnAfterItemUseFields{
		ItemId: itemId,
	}
	GameEvent.Go(OnAfterItemUse, NewEventFields(user, onAfterItemUseFields))

	err = g.afterAction(user, EffectUseInstant)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) DropItem(userId, itemId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	err = user.Inventory.DropItem(itemId)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) StartTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer.Start()
}

func (g *BaseGame) StopTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer.Stop()
}

func (g *BaseGame) GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer.GetTimeLeft(), user.Timer.IsActive(), GameSettings.NextTimerResetDate(), nil
}
