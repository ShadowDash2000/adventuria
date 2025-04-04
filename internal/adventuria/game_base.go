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

type BaseGame struct {
	gc    *GameComponents
	users *cache.MemoryCache[string, *User]
}

type GameComponents struct {
	App      core.App
	Cells    *Cells
	Log      *Log
	Cols     *collections.Collections
	Settings *Settings
	Event    Event
}

func New(app core.App) Game {
	cols := collections.NewCollections(app)
	gc := &GameComponents{
		App:   app,
		Log:   NewLog(cols, app),
		Cols:  cols,
		Event: NewEvent(),
	}

	return &BaseGame{
		gc:    gc,
		users: cache.NewMemoryCache[string, *User](0, true),
	}
}

func (g *BaseGame) Init() {
	g.gc.Settings = NewSettings(g.gc)
	g.gc.Cells = NewCells(g.gc)
}

func (g *BaseGame) Settings() *Settings {
	return g.gc.Settings
}

func (g *BaseGame) Event() Event {
	return g.gc.Event
}

func (g *BaseGame) GetUser(userId string) (*User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId, g.gc)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *BaseGame) afterAction(user *User, event string) error {
	onAfterActionFields := &OnAfterActionFields{
		Event: event,
	}
	g.Event().Go(OnAfterAction, NewEventFields(user, g.gc, onAfterActionFields))

	err := user.Save()
	if err != nil {
		return err
	}

	_, err = user.Inventory.applyEffects(event)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) ChooseGame(game string, userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != ActionTypeChooseGame {
		return errors.New("next step isn't choose game")
	}

	currentCell, _ := user.CurrentCell()

	action := NewAction(userId, ActionTypeChooseGame, g.gc)
	action.SetCell(currentCell.Id)
	action.SetValue(game)
	err = action.Save()
	if err != nil {
		return err
	}

	onAfterChooseGameFields := &OnAfterChooseGameFields{}
	g.Event().Go(OnAfterChooseGame, NewEventFields(user, g.gc, onAfterChooseGameFields))

	err = g.afterAction(user, EffectUseOnChooseGame)
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
	actionsCollection, err := g.gc.Cols.Get(TableActions)
	if err != nil {
		return err
	}

	record := &core.Record{}
	err = g.gc.App.
		RecordQuery(actionsCollection).
		AndWhere(
			dbx.HashExp{
				"user": userId,
				"id":   actionId,
			},
		).
		AndWhere(
			dbx.Or(
				dbx.HashExp{"type": ActionTypeChooseResult},
				dbx.HashExp{"type": ActionTypeDrop},
				dbx.HashExp{"type": ActionTypeReroll},
			),
		).
		Limit(1).
		One(record)
	if err != nil {
		return err
	}

	action := NewActionFromRecord(record, g.gc)
	action.SetComment(comment)
	action.SetIcon(file)

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

	action := NewAction(userId, ActionTypeReroll, g.gc)
	action.SetCell(currentCell.Id)
	action.SetComment(comment)
	action.SetValue(user.LastAction.Value())
	action.SetIcon(file)
	err = action.Save()
	if err != nil {
		return err
	}

	onAfterRerollFields := &OnAfterRerollFields{}
	g.Event().Go(OnAfterReroll, NewEventFields(user, g.gc, onAfterRerollFields))

	err = g.afterAction(user, EffectUseOnReroll)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) Roll(userId string) (int, []int, *Cell, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, nil, nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return 0, nil, nil, err
	}

	if nextStepType != ActionTypeRoll {
		return 0, nil, nil, errors.New("next step isn't roll")
	}

	onBeforeRollFields := &OnBeforeRollFields{
		Dices: []Dice{DiceTypeD6, DiceTypeD6},
	}
	g.Event().Go(OnBeforeRoll, NewEventFields(user, g.gc, onBeforeRollFields))

	onBeforeRollMoveFields := &OnBeforeRollMoveFields{
		N: 0,
	}
	diceRolls := make([]int, len(onBeforeRollFields.Dices))
	for i, dice := range onBeforeRollFields.Dices {
		diceRolls[i] = dice.Roll()
		onBeforeRollMoveFields.N += diceRolls[i]
	}

	g.Event().Go(OnBeforeRollMove, NewEventFields(user, g.gc, onBeforeRollMoveFields))

	_, currentCell, err := user.Move(onBeforeRollMoveFields.N)

	onAfterRollFields := &OnAfterRollFields{
		Dices: onBeforeRollFields.Dices,
		N:     onBeforeRollMoveFields.N,
	}
	g.Event().Go(OnAfterRoll, NewEventFields(user, g.gc, onAfterRollFields))

	err = g.afterAction(user, EffectUseOnRoll)
	if err != nil {
		return 0, nil, nil, err
	}

	return onBeforeRollMoveFields.N, diceRolls, currentCell, nil
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

	if currentCell.CantDrop() {
		return errors.New("can't drop on this cell")
	}

	onBeforeDropFields := &OnBeforeDropFields{
		IsSafeDrop: false,
	}
	g.Event().Go(OnBeforeDrop, NewEventFields(user, g.gc, onBeforeDropFields))

	action := NewAction(userId, ActionTypeDrop, g.gc)
	action.SetCell(currentCell.Id)
	action.SetComment(comment)
	action.SetValue(user.LastAction.Value())
	action.SetIcon(file)
	err = action.Save()
	if err != nil {
		return err
	}

	if !onBeforeDropFields.IsSafeDrop && !currentCell.IsSafeDrop() {
		user.SetPoints(user.Points() + g.gc.Settings.PointsForDrop())
		user.SetDropsInARow(user.DropsInARow() + 1)

		if !user.IsSafeDrop() {
			if err = user.MoveToJail(); err != nil {
				return err
			}
		}
	}

	onAfterDropFields := &OnAfterDropFields{}
	g.Event().Go(OnAfterDrop, NewEventFields(user, g.gc, onAfterDropFields))

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
	g.Event().Go(OnBeforeDone, NewEventFields(user, g.gc, onBeforeDoneFields))

	currentCell, _ := user.CurrentCell()

	action := NewAction(userId, ActionTypeChooseResult, g.gc)
	action.SetCell(currentCell.Id)
	action.SetComment(comment)
	action.SetValue(user.LastAction.Value())
	action.SetIcon(file)
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
	user.SetPoints(user.Points() + cellPoints)

	onAfterDoneFields := &OnAfterDoneFields{}
	g.Event().Go(OnAfterDone, NewEventFields(user, g.gc, onAfterDoneFields))

	err = g.afterAction(user, EffectUseOnChooseResult)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) GetLastAction(userId string) (bool, Action, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return false, nil, err
	}

	return user.IsInJail(), user.LastAction, nil
}

func (g *BaseGame) GetItemsEffects(userId, event string) (*Effects, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	effects, _, err := user.Inventory.GetEffects(event)
	if err != nil {
		return nil, err
	}

	return effects, nil
}

func (g *BaseGame) RollCell(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != ActionTypeRollCell {
		return "", errors.New("next step isn't roll cell")
	}

	gameCells := g.gc.Cells.GetAllByType(CellTypeGame)
	cell := helper.RandomItemFromSlice(gameCells)

	action := NewAction(userId, ActionTypeRollCell, g.gc)
	action.SetCell(user.LastAction.CellId())
	action.SetValue(cell.GetString("name"))
	err = action.Save()
	if err != nil {
		return "", err
	}

	onAfterWheelRollFields := &OnAfterWheelRollFields{}
	g.Event().Go(OnAfterWheelRoll, NewEventFields(user, g.gc, onAfterWheelRollFields))

	err = g.afterAction(user, "")
	if err != nil {
		return "", err
	}

	return cell.Id, nil
}

func (g *BaseGame) RollItem(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != ActionTypeRollItem {
		return "", errors.New("next step isn't roll item")
	}

	currentCell, _ := user.CurrentCell()

	item, err := RandomItem(g.gc)
	if err != nil {
		return "", err
	}

	action := NewAction(userId, ActionTypeRollItem, g.gc)
	action.SetCell(currentCell.Id)
	action.SetValue(item.Name())
	if user.ItemWheelsCount() > 0 {
		action.SetNotAffectNextStep(true)
	}
	err = action.Save()
	if err != nil {
		return "", err
	}

	err = user.Inventory.MustAddItem(item.Id)
	if err != nil {
		return "", err
	}

	// TODO: this should be in observer
	// or not... i dunno actually 😳
	if user.ItemWheelsCount() > 0 {
		user.SetItemWheelsCount(user.ItemWheelsCount() - 1)
	}

	onAfterItemRollFields := &OnAfterItemRollFields{
		ItemId: item.Id,
	}
	g.Event().Go(OnAfterItemRoll, NewEventFields(user, g.gc, onAfterItemRollFields))
	onAfterWheelRollFields := &OnAfterWheelRollFields{
		ItemId: item.Id,
	}
	g.Event().Go(OnAfterWheelRoll, NewEventFields(user, g.gc, onAfterWheelRollFields))

	err = g.afterAction(user, EffectUseOnRollItem)
	if err != nil {
		return "", err
	}

	return item.Id, nil
}

func (g *BaseGame) RollWheelPreset(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != ActionTypeRollWheelPreset {
		return "", errors.New("next step isn't roll wheel preset")
	}

	currentCell, _ := user.CurrentCell()

	item, err := RandomPresetItem(currentCell.Preset(), g.gc)
	if err != nil {
		return "", err
	}

	action := NewAction(userId, ActionTypeRollWheelPreset, g.gc)
	action.SetCell(currentCell.Id)
	action.SetValue(item.GetString("name"))
	err = action.Save()
	if err != nil {
		return "", err
	}

	onAfterWheelRollFields := &OnAfterWheelRollFields{
		ItemId: item.Id,
	}
	g.Event().Go(OnAfterWheelRoll, NewEventFields(user, g.gc, onAfterWheelRollFields))

	err = g.afterAction(user, "")
	if err != nil {
		return "", err
	}

	return item.Id, nil
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
	g.Event().Go(OnAfterItemUse, NewEventFields(user, g.gc, onAfterItemUseFields))

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

	return user.Timer.GetTimeLeft(), user.Timer.IsActive(), g.gc.Settings.NextTimerResetDate(), nil
}
