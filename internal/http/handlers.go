package http

import (
	"adventuria/internal/adventuria"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

type Handlers struct {
	Game adventuria.Game
}

func New(g adventuria.Game) *Handlers {
	return &Handlers{Game: g}
}

func (h *Handlers) NextActionTypeHandler(e *core.RequestEvent) error {
	nextStepType, err := h.Game.NextActionType(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"nextStepType": nextStepType,
	})

	return nil
}

func (h *Handlers) UpdateActionHandler(e *core.RequestEvent) error {
	data := struct {
		ActionID string `form:"actionId"`
		Comment  string `form:"comment"`
	}{}

	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err := h.Game.UpdateAction(data.ActionID, data.Comment, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollDice", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)
	return nil
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	res, err := h.Game.DoAction("reroll", e.Auth.Id, req)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)
	return nil
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	err := e.BindBody(&req)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	res, err := h.Game.DoAction("drop", e.Auth.Id, req)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)

	return nil
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	res, err := h.Game.DoAction("done", e.Auth.Id, req)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)
	return nil
}

func (h *Handlers) RollWheelHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollWheel", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)

	return nil
}

func (h *Handlers) RollItemHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollItem", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, res)

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

func (h *Handlers) GetTimeLeftByUserHandler(e *core.RequestEvent) error {
	userId := e.Request.PathValue("userId")

	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(userId)
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
