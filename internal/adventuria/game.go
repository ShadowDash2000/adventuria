package adventuria

import (
	"adventuria/pkg/cache"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"math/rand"
)

type Game struct {
	app        core.App
	users      *cache.MemoryCache[string, *User]
	cells      *cache.MemoryCache[int, *core.Record]
	cellByCode *cache.MemoryCache[string, *core.Record]
}

func NewGame(app core.App) *Game {
	return &Game{
		app:        app,
		users:      cache.NewMemoryCache[string, *User](0, true),
		cells:      cache.NewMemoryCache[int, *core.Record](0, true),
		cellByCode: cache.NewMemoryCache[string, *core.Record](0, true),
	}
}

func (g *Game) Init() error {
	err := g.fetchCells()
	if err != nil {
		return err
	}

	g.bindHooks()

	return nil
}

func (g *Game) bindHooks() {
	// CELLS
	g.app.OnRecordAfterCreateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		g.cells.Set(e.Record.GetInt("sort"), e.Record)
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			g.cellByCode.Set(cellCode, e.Record)
		}
		return e.Next()
	})
	g.app.OnRecordAfterUpdateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		g.cells.Set(e.Record.GetInt("sort"), e.Record)
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			g.cellByCode.Set(cellCode, e.Record)
		}
		return e.Next()
	})
	g.app.OnRecordAfterDeleteSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		g.cells.Delete(e.Record.GetInt("sort"))
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			g.cellByCode.Delete(cellCode)
		}
		return e.Next()
	})
}

func (g *Game) GetUser(userId string) (*User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId, g.app)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *Game) fetchCells() error {
	g.cells.Clear()
	g.cellByCode.Clear()

	cells, err := g.app.FindRecordsByFilter(
		TableCells,
		"",
		"sort",
		0,
		0,
	)
	if err != nil {
		return err
	}

	for _, cell := range cells {
		g.cells.Set(cell.GetInt("sort"), cell)
		code := cell.GetString("code")
		if code != "" {
			g.cellByCode.Set(code, cell)
		}
	}

	return nil
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

	actionsCollection, err := g.app.FindCollectionByNameOrId(TableActions)
	if err != nil {
		return nil, nil, err
	}

	cellsPassed := user.GetCellsPassed()
	currentCellNum := (cellsPassed + n) % g.cells.Count()
	currentCell, _ := g.cells.Get(currentCellNum)

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

	effects, err := user.Inventory.ApplyEffects(ItemUseTypeOnRoll)
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

	effects, err := user.Inventory.ApplyEffects(ItemUseTypeOnDrop)
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

	jailCell, ok := g.cellByCode.Get("jail")
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
	record.Set("value", cell.Id)
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

	movies, err := g.app.FindRecordsByFilter(
		TableWheelItems,
		"type = {:type}",
		"",
		0,
		0,
		dbx.Params{"type": "movie"},
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
	record.Set("value", movie.Id)
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
		return "", errors.New("items not found")
	}

	item := items[rand.Intn(len(items)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollItem)
	record.Set("value", item.Id)
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	err = user.Inventory.AddItem(item.Id)
	if err != nil {
		return "", err
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
	record.Set("value", game.Id)
	err = g.app.Save(record)
	if err != nil {
		return "", err
	}

	return game.Id, nil
}
