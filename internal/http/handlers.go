package http

import (
	"adventuria/internal/adventuria"
	"fmt"
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

	res, err := h.Game.DoAction(e.App, e.Auth.Id, "update_comment", req)
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
	res, err := h.Game.DoAction(e.App, e.Auth.Id, "rollDice", adventuria.ActionRequest{})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		})
	} else if !res.Success {
		return e.JSON(http.StatusBadRequest, res)
	}

	fmt.Println(res)

	return e.JSON(http.StatusOK, res)
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction(e.App, e.Auth.Id, "reroll", req)
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

	res, err := h.Game.DoAction(e.App, e.Auth.Id, "drop", req)
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

	res, err := h.Game.DoAction(e.App, e.Auth.Id, "done", req)
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
	res, err := h.Game.DoAction(e.App, e.Auth.Id, "rollWheel", adventuria.ActionRequest{})
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
	res, err := h.Game.DoAction(e.App, e.Auth.Id, "rollItem", adventuria.ActionRequest{})
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

func (h *Handlers) RollItemOnCellHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.App, e.Auth.Id, "rollItemOnCell", adventuria.ActionRequest{})
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

func (h *Handlers) BuyItemHandler(e *core.RequestEvent) error {
	req := adventuria.ActionRequest{}

	err := e.BindBody(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction(e.App, e.Auth.Id, "buyItem", req)
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
	data := adventuria.UseItemRequest{}
	err := e.BindBody(&data)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.Game.UseItem(e.App, e.Auth.Id, data)
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

	err = h.Game.DropItem(e.App, e.Auth.Id, data.ItemId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) StartTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StartTimer(e.App, e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) StopTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StopTimer(e.App, e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, "")
}

func (h *Handlers) GetTimeLeftHandler(e *core.RequestEvent) error {
	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(e.App, e.Auth.Id)
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

	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(e.App, userId)
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
	actions, err := h.Game.GetAvailableActions(e.App, e.Auth.Id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, actions)
}

func (h *Handlers) GetItemEffectVariants(e *core.RequestEvent) error {
	req := struct {
		InvItemId string `json:"inv_item_id"`
		EffectId  string `json:"effect_id"`
	}{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	effectVariants, err := h.Game.GetItemEffectVariants(e.App, e.Auth.Id, req.InvItemId, req.EffectId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, effectVariants)
}

func (h *Handlers) GetActionVariants(e *core.RequestEvent) error {
	action := e.Request.URL.Query().Get("action")

	actions, err := h.Game.GetActionVariants(e.App, e.Auth.Id, action)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, actions)
}
