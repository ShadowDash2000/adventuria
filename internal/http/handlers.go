package http

import (
	"adventuria/internal/adventuria"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

type Handlers struct {
	Game *adventuria.Game
}

func New(g *adventuria.Game) *Handlers {
	return &Handlers{Game: g}
}

func (h *Handlers) UpdateActionHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction("update_comment", e.Auth.Id, req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollDice", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction("reroll", e.Auth.Id, req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	err := e.BindBody(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction("drop", e.Auth.Id, req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction("done", e.Auth.Id, req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) RollWheelHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollWheel", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) RollItemHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction("rollItem", e.Auth.Id, adventuria.ActionRequest{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) UseItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.Game.UseItem(e.Auth.Id, data.ItemId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) DropItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.Game.DropItem(e.Auth.Id, data.ItemId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) StartTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StartTimer(e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) StopTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StopTimer(e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) GetTimeLeftHandler(e *core.RequestEvent) error {
	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"time":               time,
		"isActive":           isActive,
		"nextTimerResetDate": nextTimerResetDate,
	})
}

func (h *Handlers) GetTimeLeftByUserHandler(e *core.RequestEvent) error {
	userId := e.Request.PathValue("userId")

	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"time":               time,
		"isActive":           isActive,
		"nextTimerResetDate": nextTimerResetDate,
	})
}

func (h *Handlers) GetAvailableActions(e *core.RequestEvent) error {
	actions, err := h.Game.GetAvailableActions(e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, actions)
}
