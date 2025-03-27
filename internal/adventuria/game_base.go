package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"adventuria/pkg/helper"
	"errors"
	"github.com/AlexanderGrom/go-event"
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
	Event    event.Dispatcher
}

func New(app core.App) Game {
	cols := collections.NewCollections(app)
	gc := &GameComponents{
		App:   app,
		Log:   NewLog(cols, app),
		Cols:  cols,
		Event: event.New(),
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

func (g *BaseGame) Event() event.Dispatcher {
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
	err := user.Save()
	if err != nil {
		return err
	}

	g.gc.Event.Go(OnAfterAction, user, event, g.gc)

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

	g.gc.Event.Go(OnAfterChooseGame, user, g.gc)

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

	g.gc.Event.Go(OnAfterReroll, user, g.gc)

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

	dicesResult := &RollDicesResult{
		[]Dice{DiceTypeD6, DiceTypeD6},
	}
	g.gc.Event.Go(OnBeforeRoll, user, dicesResult, g.gc)

	rollResult := &RollResult{
		N: 0,
	}
	diceRolls := make([]int, len(dicesResult.Dices))
	for i, dice := range dicesResult.Dices {
		diceRolls[i] = dice.Roll()
		rollResult.N += diceRolls[i]
	}

	g.gc.Event.Go(OnBeforeRollMove, user, rollResult, g.gc)

	_, currentCell, err := user.Move(rollResult.N)

	g.gc.Event.Go(OnAfterRoll, user, rollResult, g.gc)

	err = g.afterAction(user, EffectUseOnRoll)
	if err != nil {
		return 0, nil, nil, err
	}

	return rollResult.N, diceRolls, currentCell, nil
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

	effects := &DropEffects{}
	g.gc.Event.Go(OnBeforeDrop, user, effects, g.gc)

	action := NewAction(userId, ActionTypeDrop, g.gc)
	action.SetCell(currentCell.Id)
	action.SetComment(comment)
	action.SetValue(user.LastAction.Value())
	action.SetIcon(file)
	err = action.Save()
	if err != nil {
		return err
	}

	if !effects.IsSafeDrop && !currentCell.IsSafeDrop() {
		user.SetPoints(user.Points() + g.gc.Settings.PointsForDrop())
		user.SetDropsInARow(user.DropsInARow() + 1)

		if !user.IsSafeDrop() {
			if err = user.MoveToJail(); err != nil {
				return err
			}
		}
	}

	g.gc.Event.Go(OnAfterDrop, user, g.gc)

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

	effects := &DoneEffects{}
	g.gc.Event.Go(OnBeforeDone, user, effects, g.gc)

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
	if effects.CellPointsDivide != 0 {
		cellPoints /= effects.CellPointsDivide
	}

	user.SetDropsInARow(0)
	user.SetIsInJail(false)
	user.SetPoints(user.Points() + cellPoints)

	g.gc.Event.Go(OnAfterDone, user, g.gc)

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

	g.gc.Event.Go(OnAfterWheelRoll, user, g.gc)

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

	g.gc.Event.Go(OnAfterItemRoll, user, g.gc)
	g.gc.Event.Go(OnAfterWheelRoll, user, g.gc)

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

	g.gc.Event.Go(OnAfterWheelRoll, user, g.gc)

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

	g.gc.Event.Go(OnAfterItemUse, user, g.gc)

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
