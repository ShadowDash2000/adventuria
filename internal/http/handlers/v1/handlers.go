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
	n, currentCell, err := h.Game.Roll(e.Auth)
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

	err := h.Game.ChooseGame(data.Game, e.Auth)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GetNextStepTypeHandler(e *core.RequestEvent) error {
	nextStepType, err := h.Game.GetNextStepType(e.Auth)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"nextStepType": nextStepType,
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

	err := h.Game.Reroll(data.Comment, e.Auth)
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

	err = h.Game.Drop(data.Comment, e.Auth)
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

	err := h.Game.Done(data.Comment, e.Auth)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GameResultHandler(e *core.RequestEvent) error {
	canDrop, isInJail, action, err := h.Game.GameResult(e.Auth)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	actionFields := action.FieldsData()

	e.JSON(http.StatusOK, map[string]interface{}{
		"game":     actionFields["game"].(string),
		"canDrop":  canDrop,
		"isInJail": isInJail,
	})

	return nil
}

/*func (h *Handlers) RollRandomCellHandler(e *core.RequestEvent) error {
	cellId, err := h.Game.RollRandomCell(e)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"cellId": cellId,
	})

	return nil
}*/
