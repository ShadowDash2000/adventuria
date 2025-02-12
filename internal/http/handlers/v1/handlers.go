package handlers

import (
	"adventuria/internal/usecases"
	"github.com/pocketbase/pocketbase/core"
	"net/http"
)

type Handlers struct {
	Game *usecases.Game
}

func New(g *usecases.Game) *Handlers {
	return &Handlers{Game: g}
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	n, currentCell, err := h.Game.Roll(nil, nil, e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	currentCellFields := currentCell.FieldsData()

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

func (h *Handlers) ChooseGameHandler(e *core.RequestEvent) error {
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

	err := h.Game.ChooseGame(data.Game, e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GetLastActionHandler(e *core.RequestEvent) error {
	canRoll, action, cell, err := h.Game.GetLastAction(e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	actionFields := action.FieldsData()
	cellFields := cell.FieldsData()

	e.JSON(http.StatusOK, map[string]interface{}{
		"status":   actionFields["status"].(string),
		"cellType": cellFields["type"].(string),
		"canRoll":  canRoll,
	})

	return nil
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err := h.Game.Reroll(data.Comment, e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err = h.Game.Drop(data.Comment, e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})

	return nil
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `json:"comment"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err := h.Game.Done(data.Comment, e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GameResultHandler(e *core.RequestEvent) error {
	canDrop, action, err := h.Game.GameResult(e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	actionFields := action.FieldsData()

	e.JSON(http.StatusOK, map[string]interface{}{
		"game":    actionFields["game"].(string),
		"canDrop": canDrop,
	})

	return nil
}
