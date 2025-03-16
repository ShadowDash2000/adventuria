package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"adventuria/pkg/helper"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Game struct {
	app      core.App
	log      *Log
	cols     *collections.Collections
	Settings *Settings
	cells    *Cells
	users    *cache.MemoryCache[string, *User]
}

func NewGame(app core.App) *Game {
	cols := collections.NewCollections(app)

	return &Game{
		app:   app,
		log:   NewLog(cols, app),
		cols:  cols,
		users: cache.NewMemoryCache[string, *User](0, true),
	}
}

func (g *Game) Init() {
	g.Settings = NewSettings(g.cols, g.app)
	g.cells = NewCells(g.app)
}

func (g *Game) GetUser(userId string) (*User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId, g.cells, g.Settings, g.log, g.cols, g.app)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *Game) ChooseGame(game string, userId string) error {
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

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", ActionTypeChooseGame)
	record.Set("value", game)
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseOnChooseGame)
	if err != nil {
		return err
	}

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GetNextStepType(userId string) (string, error) {
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

func (g *Game) Move(n int, userId string) (*core.Record, *core.Record, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, nil, err
	}

	actionsCollection, err := g.cols.Get(TableActions)
	if err != nil {
		return nil, nil, err
	}

	cellsPassed := user.CellsPassed()
	currentCellNum := (cellsPassed + n) % g.cells.Count()
	currentCell, _ := g.cells.GetByOrder(currentCellNum)

	record := core.NewRecord(actionsCollection)
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("value", n)
	record.Set("type", ActionTypeRoll)
	err = g.app.Save(record)
	if err != nil {
		return nil, nil, err
	}

	// Every lap gives one item wheel
	prevCellNum := cellsPassed % g.cells.Count()
	lapsPassed := (prevCellNum + n) / g.cells.Count()
	// Check if we're not moving backwards and passed new lap(-s)
	if n > 0 && lapsPassed > 0 {
		user.Set("itemWheelsCount", user.ItemWheelsCount()+lapsPassed)
	}

	user.Set("cellsPassed", cellsPassed+n)
	err = user.Save()
	if err != nil {
		return nil, nil, err
	}

	return record, currentCell, nil
}

func (g *Game) Reroll(comment string, userId string) error {
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
	cantReroll := currentCell.GetBool("cantReroll")

	if cantReroll {
		return errors.New("can't reroll on this cell")
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", ActionTypeReroll)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	user.Stats.Rerolls++
	user.Save()

	return nil
}

func (g *Game) Roll(userId string) (int, []int, *core.Record, error) {
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

	effects, err := user.Inventory.ApplyEffects(ItemUseOnRoll)
	if err != nil {
		return 0, nil, nil, err
	}

	var dices []Dice
	n := 0

	if effects.Dices != nil {
		dices = effects.Dices
	} else {
		dices = []Dice{DiceTypeD6, DiceTypeD6}
	}

	diceRolls := make([]int, len(dices))
	for i, dice := range dices {
		diceRolls[i] = dice.Roll()
		n += diceRolls[i]
	}

	if effects.DiceMultiplier > 0 {
		n *= effects.DiceMultiplier
	}
	n += effects.DiceIncrement

	if effects.RollReverse {
		n *= -1
	}

	_, currentCell, err := g.Move(n, userId)

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return n, diceRolls, currentCell, err
	}

	user.Stats.DiceRolls++
	if n > user.Stats.MaxDiceRoll {
		user.Stats.MaxDiceRoll = n
	}
	user.Save()

	return n, diceRolls, currentCell, nil
}

func (g *Game) Drop(comment string, userId string) error {
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
	cantDrop := currentCell.GetBool("cantDrop")

	if cantDrop {
		return errors.New("can't drop on this cell")
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseOnDrop)
	if err != nil {
		return err
	}

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", ActionTypeDrop)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	if !effects.IsSafeDrop && !currentCell.GetBool("isSafeDrop") {
		user.Set("points", user.Points()+g.Settings.PointsForDrop())
		user.Set("dropsInARow", user.DropsInARow()+1)

		err = user.Save()
		if err != nil {
			return err
		}

		if !user.IsSafeDrop() {
			if err = g.GoToJail(userId); err != nil {
				return err
			}
		}
	}

	user.Stats.Drops++
	user.Save()

	return nil
}

func (g *Game) Done(comment string, userId string) error {
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

	effects, err := user.Inventory.ApplyEffects(ItemUseOnChooseResult)
	if err != nil {
		return err
	}

	currentCell, _ := user.CurrentCell()

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", nextStepType)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	user.Stats.Finished++

	cellPoints := currentCell.GetInt("points")
	if effects.CellPointsDivide != 0 {
		cellPoints /= effects.CellPointsDivide
	}

	user.Set("dropsInARow", 0)
	user.Set("isInJail", false)
	user.Set("points", user.Points()+cellPoints)
	err = user.Save()
	if err != nil {
		return err
	}

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GetLastAction(userId string) (bool, *core.Record, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return false, nil, err
	}

	return user.IsInJail(), user.lastAction, nil
}

func (g *Game) GoToJail(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	jailCellPos, ok := g.cells.GetOrderByType(CellTypeJail)
	if !ok {
		return errors.New("jail cell not found")
	}

	currentCellNum := user.CellsPassed() % g.cells.Count()

	_, _, err = g.Move(jailCellPos-currentCellNum, userId)
	if err != nil {
		return err
	}

	user.Set("isInJail", true)
	user.Stats.WasInJail++
	err = user.Save()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GetItemsEffects(userId, event string) (*Effects, error) {
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

func (g *Game) RollCell(userId string) (string, error) {
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

	gameCells := g.cells.GetAllByType(CellTypeGame)
	cell := helper.RandomItemFromSlice(gameCells)

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollCell)
	record.Set("value", cell.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	user.Stats.WheelRolled++
	user.Save()

	return cell.Id, nil
}

func (g *Game) RollItem(userId string) (string, error) {
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

	items, err := g.app.FindRecordsByFilter(
		TableItems,
		"isRollable = true",
		"",
		0,
		0,
	)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "", errors.New("invItems not found")
	}

	item := helper.RandomItemFromSlice(items)

	currentCell, _ := user.CurrentCell()

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", ActionTypeRollItem)
	record.Set("value", item.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	// 'Cause we have items that actually doesn't affect the game,
	// we need to check if an item have some effects.
	// If item doesn't have any effect, we don't need to store it to inventory.
	if len(item.GetStringSlice("effects")) > 0 {
		user.Inventory.MustAddItem(item.Id)
	}

	if user.ItemWheelsCount() > 0 {
		user.Set("itemWheelsCount", user.ItemWheelsCount()-1)
		err = user.Save()
		if err != nil {
			return "", err
		}
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseOnRollItem)
	if err != nil {
		return "", err
	}

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return "", err
	}

	user.Stats.WheelRolled++
	user.Save()

	return item.Id, nil
}

func (g *Game) RollWheelPreset(userId string) (string, error) {
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

	wheelItems, err := g.app.FindRecordsByFilter(
		TableWheelItems,
		"presets.id = {:presetId}",
		"",
		0,
		0,
		dbx.Params{
			"presetId": currentCell.GetString("preset"),
		},
	)
	if err != nil {
		return "", err
	}

	if len(wheelItems) == 0 {
		return "", errors.New("wheel items for preset not found")
	}

	item := helper.RandomItemFromSlice(wheelItems)

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("type", ActionTypeRollWheelPreset)
	record.Set("value", item.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	user.Stats.WheelRolled++
	user.Save()

	return item.Id, nil
}

func (g *Game) UseItem(userId, itemId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	err = user.Inventory.UseItem(itemId)
	if err != nil {
		return err
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseInstant)
	if err != nil {
		return err
	}

	err = g.ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	user.Stats.ItemsUsed++
	user.Save()

	return nil
}

func (g *Game) ApplyGenericEffects(effects *Effects, user *User) error {
	if effects.PointsIncrement != 0 {
		user.Set("points", user.Points()+effects.PointsIncrement)
		err := user.Save()
		if err != nil {
			return err
		}
	}

	if effects.TimerIncrement != 0 {
		err := user.Timer.AddSecondsTimeLimit(effects.TimerIncrement)
		if err != nil {
			return err
		}
	}

	if effects.JailEscape {
		user.Set("isInJail", false)
		err := user.Save()
		if err != nil {
			return err
		}
	}

	if effects.DropInventory {
		err := user.Inventory.DropInventory()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) DropItem(userId, itemId string) error {
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

func (g *Game) StartTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer.Start()
}

func (g *Game) StopTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer.Stop()
}

func (g *Game) GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer.GetTimeLeft(), user.Timer.IsActive(), g.Settings.NextTimerResetDate(), nil
}
