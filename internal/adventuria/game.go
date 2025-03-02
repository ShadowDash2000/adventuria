package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"math/rand"
	"time"
)

type Game struct {
	app      core.App
	log      *Log
	cols     *collections.Collections
	settings *Settings
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
	g.settings = NewSettings(g.cols, g.app)
	g.cells = NewCells(g.app)
}

func (g *Game) GetUser(userId string) (*User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId, g.cells, g.settings, g.log, g.cols, g.app)
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

	if nextStepType != UserNextStepChooseGame {
		return errors.New("next step isn't choose game")
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeGame)
	record.Set("value", game)
	err = g.app.Save(record)
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

	cellsPassed := user.GetCellsPassed()
	currentCellNum := (cellsPassed + n) % g.cells.Count()
	currentCell, _ := g.cells.GetBySort(currentCellNum)

	record := core.NewRecord(actionsCollection)
	record.Set("user", userId)
	record.Set("cell", currentCell.Id)
	record.Set("roll", n)
	record.Set("type", ActionTypeRoll)
	err = g.app.Save(record)
	if err != nil {
		return nil, nil, err
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

	if nextStepType != UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeReroll)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

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

	if nextStepType != UserNextStepRoll {
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

	if nextStepType != UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return err
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseOnDrop)
	if err != nil {
		return err
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeDrop)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	if !effects.IsSafeDrop && cell.GetString("type") != CellTypeBigWin {
		user.Set("points", user.GetPoints()-2)
		user.Set("dropsInARow", user.GetDropsInARow()+1)

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

	if nextStepType != UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return err
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeDone)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	user.Set("dropsInARow", 0)
	user.Set("isInJail", false)
	user.Set("points", user.GetPoints()+cell.GetInt("points"))
	err = user.Save()
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

	jailCell, ok := g.cells.GetByCode("jail")
	if !ok {
		return errors.New("jail cell not found")
	}

	currentCellNum := user.GetCellsPassed() % g.cells.Count()
	jailCellPos := jailCell.GetInt("sort")

	_, _, err = g.Move(jailCellPos-currentCellNum, userId)
	if err != nil {
		return err
	}

	user.Set("isInJail", true)
	err = user.Save()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GetCellsByType(t string) []*core.Record {
	var gameCells []*core.Record
	for _, cell := range g.cells.GetAll() {
		if cell.GetString("type") == t {
			gameCells = append(gameCells, cell)
		}
	}
	return gameCells
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

	if nextStepType != UserNextStepRollJailCell {
		return "", errors.New("next step isn't roll jail cell")
	}

	gameCells := g.GetCellsByType(CellTypeGame)
	cell := gameCells[rand.Intn(len(gameCells)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollCell)
	record.Set("value", cell.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	return cell.Id, nil
}

func (g *Game) RollMovie(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != UserNextStepRollMovie {
		return "", errors.New("next step isn't roll movie")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return "", err
	}

	movies, err := g.app.FindRecordsByFilter(
		TableWheelItems,
		"type = {:type} && preset = {:preset}",
		"",
		0,
		0,
		dbx.Params{
			"type":   "movie",
			"preset": cell.GetString("preset"),
		},
	)
	if err != nil {
		return "", err
	}

	if len(movies) == 0 {
		return "", errors.New("movies not found")
	}

	movie := movies[rand.Intn(len(movies)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollMovie)
	record.Set("value", movie.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	return movie.Id, nil
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

	if nextStepType != UserNextStepRollItem {
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

	item := items[rand.Intn(len(items)-1)]

	if len(item.GetStringSlice("effects")) > 0 {
		record := core.NewRecord(user.lastAction.Collection())
		record.Set("user", userId)
		record.Set("cell", user.lastAction.GetString("cell"))
		record.Set("type", ActionTypeRollItem)
		record.Set("value", item.GetString("name"))
		err = g.app.Save(record)
		if err != nil {
			return "", err
		}

		err = user.Inventory.AddItem(item.Id)
		if err != nil {
			return "", err
		}
	}

	effects, err := user.Inventory.ApplyEffects(ItemUseOnRollItem)
	if err != nil {
		return "", err
	}

	if effects.DropInventory {
		err = user.Inventory.DropInventory()
		if err != nil {
			return "", err
		}
	}

	return item.Id, nil
}

func (g *Game) RollBigWin(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != UserNextStepRollBigWin {
		return "", errors.New("next step isn't roll movie")
	}

	games, err := g.app.FindRecordsByFilter(
		TableWheelItems,
		"type = {:type}",
		"",
		0,
		0,
		dbx.Params{"type": "legendaryGame"},
	)
	if err != nil {
		return "", err
	}

	if len(games) == 0 {
		return "", errors.New("legendary games not found")
	}

	game := games[rand.Intn(len(games)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollBigWin)
	record.Set("value", game.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	return game.Id, nil
}

func (g *Game) RollDeveloper(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	if nextStepType != UserNextStepRollDeveloper {
		return "", errors.New("next step isn't roll developer")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return "", err
	}

	games, err := g.app.FindRecordsByFilter(
		TableWheelItems,
		"type = {:type} && preset = {:preset}",
		"",
		0,
		0,
		dbx.Params{
			"type":   "developer",
			"preset": cell.GetString("preset"),
		},
	)
	if err != nil {
		return "", err
	}

	if len(games) == 0 {
		return "", errors.New("games not found")
	}

	game := games[rand.Intn(len(games)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollDeveloper)
	record.Set("value", game.GetString("name"))
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	return game.Id, nil
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

	if effects.TimerIncrement != 0 {
		err = user.Timer.AddSecondsTimeLimit(effects.TimerIncrement)
		if err != nil {
			return err
		}
	}

	if effects.JailEscape {
		user.Set("isInJail", false)
		err = user.Save()
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

func (g *Game) MovieDone(comment string, userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != UserNextStepMovieResult {
		return errors.New("next step isn't choose movie result")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return err
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeMovieResult)
	record.Set("comment", comment)
	record.Set("value", user.lastAction.GetString("value"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	user.Set("dropsInARow", 0)
	user.Set("isInJail", false)
	user.Set("points", user.GetPoints()+cell.GetInt("points"))
	err = user.Save()
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

func (g *Game) GetTimeLeft(userId string) (time.Duration, bool, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, false, err
	}

	return user.Timer.GetTimeLeft(), user.Timer.IsActive(), nil
}
