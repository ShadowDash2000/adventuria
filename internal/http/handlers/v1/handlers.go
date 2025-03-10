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

	e.JSON(http.StatusOK, map[string]interface{}{
		"roll":      n,
		"diceRolls": diceRolls,
		"cellId":    currentCell.Id,
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

func (h *Handlers) GetLastActionHandler(e *core.RequestEvent) error {
	isInJail, action, err := h.Game.GetLastAction(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"title":    action.GetString("value"),
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

func (h *Handlers) RollWheelPresetHandler(e *core.RequestEvent) error {
	itemId, err := h.Game.RollWheelPreset(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"itemId": itemId,
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
	effects, err := h.Game.GetItemsEffects(e.Auth.Id, adventuria.ItemUseOnRoll)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, effects)

	return nil
}

func (h *Handlers) UseItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err = h.Game.UseItem(e.Auth.Id, data.ItemId)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) DropItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err = h.Game.DropItem(e.Auth.Id, data.ItemId)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) StartTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StartTimer(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) StopTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StopTimer(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) GetTimeLeftHandler(e *core.RequestEvent) error {
	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"time":               time,
		"isActive":           isActive,
		"nextTimerResetDate": nextTimerResetDate,
	})

	return nil
}
