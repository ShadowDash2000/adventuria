package usecases

import (
	"adventuria/internal/adventuria"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"math/rand/v2"
	"net/http"
	"slices"
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) ChooseGame(game string, e *core.RequestEvent) error {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	if len(actions) == 0 {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action := actions[0]
	errs := e.App.ExpandRecord(action, []string{"cell"}, nil)
	if len(errs) > 0 {
		e.JSON(http.StatusInternalServerError, errs)
		return nil
	}

	cell := action.ExpandedOne("cell")
	cellFields := cell.FieldsData()
	actionFields := action.FieldsData()

	if cellFields["type"].(string) == "special" {
		e.JSON(http.StatusBadRequest, "You must choose a special position")
		return nil
	}

	if cellFields["code"].(string) == adventuria.SpecialCellBigWin &&
		actionFields["status"].(string) != adventuria.ActionStatusReroll {
		e.JSON(http.StatusNotFound, "You can choose a game on Big Win only if last one was rerolled")
		return nil
	}

	statuses := []string{
		adventuria.ActionStatusNotChosen,
		adventuria.ActionStatusReroll,
		adventuria.ActionStatusDrop,
	}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	statuses = []string{
		adventuria.ActionStatusReroll,
		adventuria.ActionStatusDrop,
	}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		record := core.NewRecord(action.Collection())
		record.Set("user", e.Auth.Id)
		record.Set("cell", actionFields["cell"].(string))
		record.Set("status", adventuria.ActionStatusInProgress)
		record.Set("game", game)
		err = e.App.Save(record)
		if err != nil {
			e.JSON(http.StatusInternalServerError, err.Error())
			return err
		}
	} else {
		action.Set("status", adventuria.ActionStatusInProgress)
		action.Set("game", game)
		err = e.App.Save(action)
		if err != nil {
			e.JSON(http.StatusInternalServerError, err.Error())
			return err
		}
	}

	return nil
}

func (g *Game) GetLastAction(e *core.RequestEvent) (bool, *core.Record, *core.Record, error) {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		return false, nil, nil, err
	}

	if len(actions) == 0 {
		e.JSON(http.StatusOK, map[string]interface{}{
			"status":  "",
			"canRoll": true,
		})
		return true, nil, nil, nil
	}

	action := actions[0]
	actionFields := action.FieldsData()
	errs := e.App.ExpandRecord(action, []string{"cell"}, nil)
	if len(errs) > 0 {
		for _, err = range errs {
			return false, nil, nil, err
		}
	}

	cell := action.ExpandedOne("cell")
	cellFields := cell.FieldsData()
	canRoll := true

	statuses := []string{
		adventuria.ActionStatusNotChosen,
		adventuria.ActionStatusReroll,
		adventuria.ActionStatusDrop,
		adventuria.ActionStatusInProgress,
	}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		canRoll = false
	}

	if cellFields["code"].(string) == adventuria.SpecialCellBigWin &&
		actionFields["status"].(string) == adventuria.ActionStatusReroll {
		canRoll = true
	}

	return canRoll, action, cell, err
}

func (g *Game) Reroll(comment string, e *core.RequestEvent) error {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	if len(actions) == 0 {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action := actions[0]
	actionFields := action.FieldsData()
	statuses := []string{adventuria.ActionStatusInProgress}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action.Set("status", adventuria.ActionStatusReroll)
	action.Set("comment", comment)
	err = e.App.Save(action)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func (g *Game) Move(n *int, status *string, e *core.RequestEvent) (*core.Record, *core.Record, error) {
	user, err := e.App.FindRecordById("users", e.Auth.Id)
	if err != nil {
		return nil, nil, err
	}

	actionsCollection, err := e.App.FindCollectionByNameOrId("actions")
	if err != nil {
		return nil, nil, err
	}

	cells, err := e.App.FindRecordsByFilter(
		"cells",
		"",
		"sort",
		-1,
		-1,
	)
	if err != nil {
		return nil, nil, err
	}

	if status == nil {
		status = new(string)
		*status = adventuria.ActionStatusNotChosen
	}

	userFields := user.FieldsData()
	cellsPassed := int(userFields["cellsPassed"].(float64))
	currentCellNum := (cellsPassed + *n) % len(cells)
	currentCell := cells[currentCellNum]

	record := core.NewRecord(actionsCollection)
	record.Set("user", e.Auth.Id)
	record.Set("cell", currentCell.Id)
	record.Set("roll", n)
	record.Set("status", *status)
	err = e.App.Save(record)
	if err != nil {
		return nil, nil, err
	}

	user.Set("cellsPassed", cellsPassed+*n)
	err = e.App.Save(user)
	if err != nil {
		return nil, nil, err
	}

	return record, currentCell, nil
}

func (g *Game) Roll(n *int, status *string, e *core.RequestEvent) (int, *core.Record, error) {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		return 0, nil, err
	}

	if len(actions) > 0 {
		action := actions[0]
		actionFields := action.FieldsData()
		statuses := []string{adventuria.ActionStatusDone}

		if actionFields["status"] != "" && !slices.Contains(statuses, actionFields["status"].(string)) {
			return 0, nil, errors.New("last action isn't done yet")
		}
	}

	if n == nil {
		n = new(int)
		*n = rand.IntN(6-1) + 1
	}

	_, currentCell, err := g.Move(n, status, e)

	return *n, currentCell, nil
}

func (g *Game) CanDrop(e *core.RequestEvent) (bool, error) {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		3,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		return false, err
	}

	if len(actions) == 0 {
		return false, nil
	}

	if len(actions) == 3 {
		previousActions := actions[1:3]
		i := 0

		for _, previousAction := range previousActions {
			previousActionFields := previousAction.FieldsData()

			if previousActionFields["status"].(string) == adventuria.ActionStatusDrop {
				i++
			}
		}

		if i >= 2 {
			return false, nil
		}
	}

	action := actions[0]
	actionFields := action.FieldsData()

	if actionFields["status"].(string) != adventuria.ActionStatusInProgress {
		return false, nil
	}

	return true, nil
}

func (g *Game) Drop(comment string, e *core.RequestEvent) error {
	canDrop, err := g.CanDrop(e)
	if err != nil {
		return err
	}

	if !canDrop {
		return errors.New("not allowed to drop")
	}

	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		2,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		return err
	}

	action := actions[0]
	actionFields := action.FieldsData()

	errs := e.App.ExpandRecord(action, []string{"user", "cell"}, nil)
	if len(errs) > 0 {
		for _, err = range errs {
			return err
		}
	}

	user := action.ExpandedOne("user")
	cell := action.ExpandedOne("cell")
	userFields := user.FieldsData()
	cellFields := cell.FieldsData()

	if cellFields["code"].(string) != adventuria.SpecialCellBigWin {
		points := userFields["points"].(float64) - 2

		user.Set("points", points)

		err = e.App.Save(user)
		if err != nil {
			return err
		}
	}

	action.Set("status", adventuria.ActionStatusDrop)
	action.Set("comment", comment)
	err = e.App.Save(action)
	if err != nil {
		return err
	}

	if len(actions) > 1 {
		previousAction := actions[1]
		previousActionFields := previousAction.FieldsData()

		if previousActionFields["status"].(string) == adventuria.ActionStatusDrop {
			cells, err := e.App.FindRecordsByFilter(
				"cells",
				"",
				"sort",
				-1,
				-1,
			)
			if err != nil {
				return err
			}

			if len(cells) == 0 {
				return errors.New("no cells found")
			}

			var jailCell *core.Record
			for _, cell := range cells {
				cellFields := cell.FieldsData()
				if cellFields["code"].(string) == adventuria.SpecialCellJail {
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
			roll := jailCellPos - currentCellNum

			status := ""

			_, _, err = g.Move(&roll, &status, e)
			if err != nil {
				return err
			}
		} else if actionFields["status"].(string) == adventuria.ActionStatusInProgress {
			record := core.NewRecord(action.Collection())
			record.Set("user", e.Auth.Id)
			record.Set("cell", actionFields["cell"].(string))
			record.Set("status", adventuria.ActionStatusNotChosen)
			err = e.App.Save(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Game) Done(comment string, e *core.RequestEvent) error {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	if len(actions) == 0 {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action := actions[0]
	actionFields := action.FieldsData()
	statuses := []string{adventuria.ActionStatusInProgress}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action.Set("status", adventuria.ActionStatusDone)
	action.Set("comment", comment)
	err = e.App.Save(action)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	errs := e.App.ExpandRecord(action, []string{"user"}, nil)
	if len(errs) > 0 {
		e.JSON(http.StatusInternalServerError, errs)
		return nil
	}

	errs = e.App.ExpandRecord(action, []string{"cell"}, nil)
	if len(errs) > 0 {
		e.JSON(http.StatusInternalServerError, errs)
		return nil
	}

	user := action.ExpandedOne("user")
	userFields := user.FieldsData()
	cell := action.ExpandedOne("cell")
	cellFields := cell.FieldsData()

	user.Set("points", userFields["points"].(float64)+cellFields["points"].(float64))
	err = e.App.Save(user)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func (g *Game) GameResult(e *core.RequestEvent) (bool, *core.Record, error) {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		2,
		0,
		dbx.Params{"userId": e.Auth.Id},
	)
	if err != nil {
		return false, nil, err
	}

	if len(actions) == 0 {
		return false, nil, errors.New("no active actions to record game result")
	}

	canDrop, err := g.CanDrop(e)
	if err != nil {
		return false, nil, err
	}

	return canDrop, actions[0], nil
}
