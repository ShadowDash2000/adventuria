package usecases

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/cache"
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"slices"
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
	g.app.OnRecordAfterCreateSuccess(adventuria.TableCells).BindFunc(func(e *core.RecordEvent) error {
		g.cells.Set(e.Record.GetInt("sort"), e.Record)
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			g.cellByCode.Set(cellCode, e.Record)
		}
		return e.Next()
	})
	g.app.OnRecordAfterDeleteSuccess(adventuria.TableCells).BindFunc(func(e *core.RecordEvent) error {
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
		adventuria.TableCells,
		"",
		"sort",
		0,
		0,
	)
	if err != nil {
		return err
	}

	if len(cells) == 0 {
		return errors.New("no cells found")
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

func (g *Game) ChooseGame(game string, auth *core.Record) error {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != adventuria.UserNextStepChooseGame {
		return errors.New("next step isn't choose game")
	}

	action := user.actions[0]
	actionFields := action.FieldsData()

	statuses := []string{
		adventuria.ActionStatusReroll,
		adventuria.ActionStatusDrop,
	}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		record := core.NewRecord(action.Collection())
		record.Set("user", auth.Id)
		record.Set("cell", actionFields["cell"].(string))
		record.Set("status", adventuria.ActionStatusInProgress)
		record.Set("game", game)
		err = g.app.Save(record)
		if err != nil {
			return err
		}
	} else {
		action.Set("status", adventuria.ActionStatusInProgress)
		action.Set("game", game)
		err = g.app.Save(action)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) GetNextStepType(auth *core.Record) (string, error) {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return "", err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return "", err
	}

	return nextStepType, nil
}

func (g *Game) Move(n int, status string, auth *core.Record) (*core.Record, *core.Record, error) {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return nil, nil, err
	}

	actionsCollection, err := g.app.FindCollectionByNameOrId(adventuria.TableActions)
	if err != nil {
		return nil, nil, err
	}

	if status == "" {
		status = adventuria.ActionStatusNone
	}

	cellsPassed := user.GetCellsPassed()
	currentCellNum := (cellsPassed + n) % g.cells.Count()
	currentCell, _ := g.cells.Get(currentCellNum)

	record := core.NewRecord(actionsCollection)
	record.Set("user", auth.Id)
	record.Set("cell", currentCell.Id)
	record.Set("roll", n)
	record.Set("status", status)
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

func (g *Game) Reroll(comment string, auth *core.Record) error {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != adventuria.UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	action := user.actions[0]
	action.Set("status", adventuria.ActionStatusReroll)
	action.Set("comment", comment)
	err = g.app.Save(action)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Roll(auth *core.Record) (int, *core.Record, error) {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return 0, nil, err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return 0, nil, err
	}

	if nextStepType != adventuria.UserNextStepRoll {
		return 0, nil, errors.New("next step isn't roll")
	}

	n := DiceTypeD4.Roll()

	_, currentCell, err := g.Move(n, adventuria.ActionStatusGameNotChosen, auth)

	return n, currentCell, nil
}

func (g *Game) Drop(comment string, auth *core.Record) error {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != adventuria.UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return err
	}

	action := user.actions[0]
	actionFields := action.FieldsData()

	if cell.GetString("type") != adventuria.CellTypeBigWin {
		points := user.GetPoints() - 2

		user.Set("points", points)

		err = user.Save()
		if err != nil {
			return err
		}
	}

	action.Set("status", adventuria.ActionStatusDrop)
	action.Set("comment", comment)
	err = g.app.Save(action)
	if err != nil {
		return err
	}

	if len(user.actions) > 1 {
		previousAction := user.actions[1]
		previousActionFields := previousAction.FieldsData()

		if previousActionFields["status"].(string) == adventuria.ActionStatusDrop {
			jailCell, ok := g.cellByCode.Get("jail")
			if !ok {
				return errors.New("jail cell not found")
			}

			currentCellNum := user.GetCellsPassed() % g.cells.Count()
			jailCellPos := jailCell.GetInt("sort")

			_, _, err = g.Move(
				jailCellPos-currentCellNum,
				adventuria.ActionStatusGameNotChosen,
				auth,
			)
			if err != nil {
				return err
			}
		} else {
			record := core.NewRecord(action.Collection())
			record.Set("user", auth.Id)
			record.Set("cell", actionFields["cell"].(string))
			record.Set("status", adventuria.ActionStatusGameNotChosen)
			err = g.app.Save(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Game) Done(comment string, auth *core.Record) error {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return err
	}

	nextStepType, err := user.GetNextStepType()
	if err != nil {
		return err
	}

	if nextStepType != adventuria.UserNextStepChooseResult {
		return errors.New("next step isn't choose result")
	}

	cell, err := user.GetCurrentCell()
	if err != nil {
		return err
	}

	action := user.actions[0]
	action.Set("status", adventuria.ActionStatusDone)
	action.Set("comment", comment)
	err = g.app.Save(action)
	if err != nil {
		return err
	}

	user.Set("points", user.GetPoints()+cell.GetInt("points"))
	err = user.Save()
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GameResult(auth *core.Record) (bool, bool, *core.Record, error) {
	user, err := g.GetUser(auth.Id)
	if err != nil {
		return false, false, nil, err
	}

	if len(user.actions) == 0 {
		return false, false, nil, errors.New("no active actions to record game result")
	}

	canDrop, err := user.CanDrop()
	if err != nil {
		return false, false, nil, err
	}

	isInJail, err := user.IsInJail()
	if err != nil {
		return false, false, nil, err
	}

	return canDrop, isInJail, user.actions[0], nil
}

// RollRandomCell
// Роллить рандомную клетку можно только находясь в тюрьме.
// Можно роллить только из клеток типа - game.
/*func (g *Game) RollRandomCell(e *core.RequestEvent) (string, error) {
	isInJail, err := g.user.IsInJail()
	if err != nil {
		return "", err
	}

	if !isInJail {
		return "", errors.New("you can roll random cell only in jail")
	}

	cells, err := e.App.FindRecordsByFilter(
		"cells",
		"type = game",
		"sort",
		-1,
		-1,
	)
	if err != nil {
		return "", err
	}

	n := rand.IntN(len(cells)-1) + 1

	cell := cells[n]

	return cell.Id, nil
}*/
