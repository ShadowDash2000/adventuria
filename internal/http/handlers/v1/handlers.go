package handlers

import (
	"adventuria/internal/adventuria"
	"github.com/pocketbase/pocketbase/core"
	"net/http"
)

type Handlers struct {
	Game *adventuria.Game
}

func New(g *adventuria.Game) *Handlers {
	return &Handlers{Game: g}
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	n, diceRolls, currentCell, err := h.Game.Roll(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	currentCellFields := currentCell.FieldsData()

	e.JSON(http.StatusOK, map[string]interface{}{
		"roll":      n,
		"diceRolls": diceRolls,
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

	err := h.Game.ChooseGame(data.Game, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GetNextStepTypeHandler(e *core.RequestEvent) error {
	nextStepType, err := h.Game.GetNextStepType(e.Auth.Id)
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

	err := h.Game.Reroll(data.Comment, e.Auth.Id)
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

	err = h.Game.Drop(data.Comment, e.Auth.Id)
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

	err := h.Game.Done(data.Comment, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GameResultHandler(e *core.RequestEvent) error {
	isInJail, action, err := h.Game.GetLastAction(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"game":     action.GetString("value"),
		"isInJail": isInJail,
	})

	return nil
}

func (h *Handlers) RollCellHandler(e *core.RequestEvent) error {
	cellId, err := h.Game.RollCell(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"itemId": cellId,
	})

	return nil
}

func (h *Handlers) RollMovieHandler(e *core.RequestEvent) error {
	movieId, err := h.Game.RollMovie(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"itemId": movieId,
	})

	return nil
}

func (h *Handlers) RollItemHandler(e *core.RequestEvent) error {
	itemId, err := h.Game.RollItem(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"itemId": itemId,
	})

	return nil
}

func (h *Handlers) GetRollEffectsHandler(e *core.RequestEvent) error {
	effects, err := h.Game.GetItemsEffects(e.Auth.Id, adventuria.ItemUseTypeOnRoll)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, effects)

	return nil
}
