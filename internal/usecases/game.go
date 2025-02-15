package usecases

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/cache"
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"slices"
	"time"
)

type Game struct {
	app   core.App
	users *cache.MemoryCache[*User]
	cells []*core.Record
}

func NewGame(app core.App) *Game {
	return &Game{
		app:   app,
		users: cache.NewMemoryCache[*User](30*time.Minute, false),
	}
}

func (g *Game) GetUser(auth *core.Record) (*User, error) {
	user, ok := g.users.Get(auth.Id)
	if ok {
		return user, nil
	}

	user, err := NewUser(g.app, auth)
	if err != nil {
		return nil, err
	}

	g.users.Set(auth.Id, user)
	return user, nil
}

func (g *Game) GetCells() ([]*core.Record, error) {
	if len(g.cells) > 0 {
		return g.cells, nil
	}

	cells, err := g.app.FindRecordsByFilter(
		"cells",
		"",
		"sort",
		-1,
		-1,
	)
	if err != nil {
		return nil, err
	}

	if len(cells) == 0 {
		return nil, errors.New("no cells found")
	}

	g.cells = cells

	return g.cells, nil
}

func (g *Game) ChooseGame(game string, auth *core.Record) error {
	user, err := g.GetUser(auth)
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
		err = user.AddAction(record)
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
	user, err := g.GetUser(auth)
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
	user, err := g.GetUser(auth)
	if err != nil {
		return nil, nil, err
	}

	actionsCollection, err := g.app.FindCollectionByNameOrId(adventuria.TableActions)
	if err != nil {
		return nil, nil, err
	}

	cells, err := g.GetCells()
	if err != nil {
		return nil, nil, err
	}

	if status == "" {
		status = adventuria.ActionStatusNone
	}

	userFields := user.user.FieldsData()
	cellsPassed := int(userFields["cellsPassed"].(float64))
	currentCellNum := (cellsPassed + n) % len(cells)
	currentCell := cells[currentCellNum]

	record := core.NewRecord(actionsCollection)
	record.Set("user", auth.Id)
	record.Set("cell", currentCell.Id)
	record.Set("roll", n)
	record.Set("status", status)
	err = user.AddAction(record)
	if err != nil {
		return nil, nil, err
	}

	user.user.Set("cellsPassed", cellsPassed+n)
	err = g.app.Save(user.user)
	if err != nil {
		return nil, nil, err
	}

	return record, currentCell, nil
}

func (g *Game) Reroll(comment string, auth *core.Record) error {
	user, err := g.GetUser(auth)
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
	user, err := g.GetUser(auth)
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
	user, err := g.GetUser(auth)
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
	userFields := user.user.FieldsData()
	cellFields := cell.FieldsData()

	if cellFields["code"].(string) != adventuria.CellTypeBigWin {
		points := userFields["points"].(float64) - 2

		user.user.Set("points", points)

		err = g.app.Save(user.user)
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
			cells, err := g.GetCells()
			if err != nil {
				return err
			}

			if len(cells) == 0 {
				return errors.New("no cells found")
			}

			var jailCell *core.Record
			for _, cell := range cells {
				cellFields := cell.FieldsData()
				if cellFields["code"].(string) == adventuria.CellTypeJail {
					jailCell = cell
					break
				}
			}

			if jailCell == nil {
				return errors.New("jail cell not found")
			}

			jailCellFields := jailCell.FieldsData()

			cellsPassed := int(userFields["cellsPassed"].(float64))
			currentCellNum := cellsPassed % len(cells)

			jailCellPos := int(jailCellFields["sort"].(float64))

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
			err = user.AddAction(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Game) Done(comment string, auth *core.Record) error {
	user, err := g.GetUser(auth)
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

	userFields := user.user.FieldsData()
	cellFields := cell.FieldsData()

	user.user.Set("points", userFields["points"].(float64)+cellFields["points"].(float64))
	err = g.app.Save(user.user)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) GameResult(auth *core.Record) (bool, bool, *core.Record, error) {
	user, err := g.GetUser(auth)
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
