package adventuria

import (
	"adventuria/pkg/cache"
	"errors"
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
	record.Set("game", game)
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
	record.Set("game", user.lastAction.GetString("game"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Roll(userId string) (int, *core.Record, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return 0, nil, err
	}

	if nextStepType != UserNextStepRoll {
		return 0, nil, errors.New("next step isn't roll")
	}

	effects, err := user.Inventory.ApplyOnRollEffects()
	if err != nil {
		return 0, nil, err
	}

	var dices []Dice
	n := 0

	if effects.Dices != nil {
		dices = effects.Dices
	} else {
		dices = []Dice{DiceTypeD6, DiceTypeD6}
	}

	for _, dice := range dices {
		n += dice.Roll()
	}

	n *= effects.DiceMultiplier
	n += effects.DiceIncrement

	_, currentCell, err := g.Move(n, userId)

	return n, currentCell, nil
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

	effects, err := user.Inventory.ApplyOnDropEffects()
	if err != nil {
		return err
	}

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeDrop)
	record.Set("comment", comment)
	record.Set("game", user.lastAction.GetString("game"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

	if !effects.IsSafeDrop && cell.GetString("type") != CellTypeBigWin {
		points := user.GetPoints() - 2

		user.Set("points", points)

		err = user.Save()
		if err != nil {
			return err
		}

		var isSafeDrop bool
		isSafeDrop, err = user.IsSafeDrop()
		if err != nil {
			return err
		}

		if !isSafeDrop {
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
	record.Set("game", user.lastAction.GetString("game"))
	err = g.app.Save(record)
	if err != nil {
		return err
	}

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
	g.cells.Range(func(k, v interface{}) bool {
		cell := v.(*core.Record)
		if cell.GetString("type") == t {
			gameCells = append(gameCells, cell)
		}
		return true
	})

	return gameCells
}

func (g *Game) RollCell(userId string) (*core.Record, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return nil, err
	}

	if nextStepType != UserNextStepRollJailCell {
		return nil, errors.New("next step isn't roll jail cell")
	}

	gameCells := g.GetCellsByType(CellTypeGame)
	cell := gameCells[rand.Intn(len(gameCells)-1)]

	record := core.NewRecord(user.lastAction.Collection())
	record.Set("user", userId)
	record.Set("cell", user.lastAction.GetString("cell"))
	record.Set("type", ActionTypeRollCell)
	record.Set("game", user.lastAction.GetString("game"))
	err = g.app.Save(record)
	if err != nil {
		return nil, err
	}

	return cell, nil
}
