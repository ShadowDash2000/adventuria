package handlers

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"math/rand/v2"
	"net/http"
	"slices"
)

const (
	ActionStatusNotChosen  string = "notChosen"
	ActionStatusReroll     string = "reroll"
	ActionStatusDrop       string = "drop"
	ActionStatusDone       string = "done"
	ActionStatusInProgress string = "inProgress"

	SpecialCellStart  string = "start"
	SpecialCellJail   string = "jail"
	SpecialCellBigWin string = "big-win"
	SpecialCellPreset string = "preset"
)

func RollHandler(e *core.RequestEvent) error {
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

	if len(actions) > 0 {
		action := actions[0]
		actionFields := action.FieldsData()
		statuses := []string{ActionStatusDone}

		if actionFields["status"] != "" && !slices.Contains(statuses, actionFields["status"].(string)) {
			e.JSON(http.StatusNotFound, "You must complete the last action")
			return nil
		}
	}

	user, err := e.App.FindRecordById("users", e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	userFields := user.FieldsData()

	cells, err := e.App.FindRecordsByFilter(
		"cells",
		"",
		"sort",
		-1,
		-1,
	)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	n := rand.IntN(6-1) + 1

	cellsPassed := int(userFields["cellsPassed"].(float64))
	currentCellNum := (cellsPassed + n) % len(cells)
	currentCell := cells[currentCellNum]
	currentCellFields := currentCell.FieldsData()

	actionsCollection, err := e.App.FindCollectionByNameOrId("actions")
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	record := core.NewRecord(actionsCollection)
	record.Set("user", e.Auth.Id)
	record.Set("cell", currentCell.Id)
	record.Set("roll", n)
	record.Set("status", ActionStatusNotChosen)
	err = e.App.Save(record)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	user.Set("cellsPassed", cellsPassed+n)
	err = e.App.Save(user)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"roll": n,
		"cell": map[string]interface{}{
			"name":        currentCellFields["name"].(string),
			"description": currentCellFields["description"].(string),
			"icon": "/api/files/" +
				currentCell.Collection().Id +
				"/" + currentCell.Id + "/" +
				currentCellFields["icon"].(string),
		},
	})
	return nil
}

func ChooseGameHandler(e *core.RequestEvent) error {
	data := struct {
		Game string `json:"game"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	if data.Game == "" {
		e.JSON(http.StatusBadRequest, "You must choose a game")
		return nil
	}

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

	if cellFields["code"].(string) == SpecialCellBigWin && actionFields["status"].(string) != ActionStatusReroll {
		e.JSON(http.StatusNotFound, "You can choose a game on Big Win only if last one was rerolled")
		return nil
	}

	statuses := []string{ActionStatusNotChosen, ActionStatusReroll, ActionStatusDrop}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	statuses = []string{ActionStatusReroll, ActionStatusDrop}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		record := core.NewRecord(action.Collection())
		record.Set("user", e.Auth.Id)
		record.Set("cell", actionFields["cell"].(string))
		record.Set("status", ActionStatusInProgress)
		record.Set("game", data.Game)
		err = e.App.Save(record)
		if err != nil {
			e.JSON(http.StatusInternalServerError, err.Error())
			return err
		}
	} else {
		action.Set("status", ActionStatusInProgress)
		action.Set("game", data.Game)
		err = e.App.Save(action)
		if err != nil {
			e.JSON(http.StatusInternalServerError, err.Error())
			return err
		}
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func GetLastActionHandler(e *core.RequestEvent) error {
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
		e.JSON(http.StatusOK, map[string]interface{}{
			"status":  "",
			"canRoll": true,
		})
		return nil
	}

	action := actions[0]
	actionFields := action.FieldsData()
	errs := e.App.ExpandRecord(action, []string{"cell"}, nil)
	if len(errs) > 0 {
		e.JSON(http.StatusInternalServerError, errs)
		return nil
	}

	cell := action.ExpandedOne("cell")
	cellFields := cell.FieldsData()
	canRoll := true

	statuses := []string{ActionStatusNotChosen, ActionStatusReroll, ActionStatusDrop, ActionStatusInProgress}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		canRoll = false
	}

	if cellFields["code"].(string) == SpecialCellBigWin && actionFields["status"].(string) == ActionStatusReroll {
		canRoll = true
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"status":   actionFields["status"].(string),
		"cellType": cellFields["type"].(string),
		"canRoll":  canRoll,
	})

	return nil
}

func RerollHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

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
	statuses := []string{ActionStatusInProgress}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action.Set("status", ActionStatusReroll)
	action.Set("comment", data.Comment)
	err = e.App.Save(action)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})

	return nil
}

func DropHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
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
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	if len(actions) == 0 {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action := actions[0]
	actionFields := action.FieldsData()

	if len(actions) > 1 {
		previousAction := actions[1]
		previousActionFields := previousAction.FieldsData()

		if previousActionFields["status"].(string) == ActionStatusDrop {
			e.JSON(http.StatusForbidden, "You must complete the last action")
			return nil
		}
	}

	if actionFields["status"].(string) != ActionStatusInProgress {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action.Set("status", ActionStatusDrop)
	action.Set("comment", data.Comment)
	err = e.App.Save(action)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})

	return nil
}

func DoneHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

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
	statuses := []string{ActionStatusInProgress}

	if !slices.Contains(statuses, actionFields["status"].(string)) {
		e.JSON(http.StatusNotFound, "You must complete the last action")
		return nil
	}

	action.Set("status", ActionStatusDone)
	action.Set("comment", data.Comment)
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

	e.JSON(http.StatusOK, struct{}{})

	return nil
}

func GameResultHandler(e *core.RequestEvent) error {
	actions, err := e.App.FindRecordsByFilter(
		"actions",
		"user.id = {:userId}",
		"-created",
		2,
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

	canDrop := true

	if len(actions) > 1 {
		previousAction := actions[1]
		previousActionFields := previousAction.FieldsData()

		if previousActionFields["status"].(string) == ActionStatusDrop {
			canDrop = false
		}
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"game":    actionFields["game"].(string),
		"canDrop": canDrop,
	})

	return nil
}
