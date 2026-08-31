package http

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/actions/custom/buy"
	"adventuria/internal/adventuria/actions/custom/done"
	"adventuria/internal/adventuria/actions/custom/drop"
	"adventuria/internal/adventuria/actions/custom/reroll"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/http/response"
	"context"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

type game interface {
	DoAction(ctx context.Context, pb core.App, playerId string, actionType model.ActionType, req model.ActionRequest) (any, error)
	UseItem(ctx context.Context, pb core.App, playerId string, itemId string, data map[string]any) error
	DropItem(ctx context.Context, pb core.App, playerId, itemId string) error
	GetAvailableActions(ctx context.Context, playerId string) ([]model.ActionType, error)
	GetEffectView(ctx context.Context, playerId, effectId string) (any, error)
	GetActionView(ctx context.Context, playerId string, actionType model.ActionType) (any, error)
	IsActionsBlocked(ctx context.Context) error
}

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
	IsEventEnded(ctx context.Context) (bool, error)
}

type Handlers struct {
	game     game
	settings settings
}

func New(game game, settings settings) *Handlers {
	return &Handlers{
		game:     game,
		settings: settings,
	}
}

func (h *Handlers) StartHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeStart, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollDice, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	req := reroll.Request{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeReroll, req)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	req := drop.Request{}

	err := e.BindBody(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeDrop, req)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	req := done.Request{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeDone, req)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) GenerateWheelHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeGenerateWheel, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RollWheelHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollWheel, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RollItemHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollItem, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RollItemOnCellHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollItemOnCell, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) BuyItemHandler(e *core.RequestEvent) error {
	req := buy.Request{}

	err := e.BindBody(&req)
	if err != nil {
		return response.Error(e, err)
	}

	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeBuy, req)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) UseItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string         `json:"item_id"`
		Data   map[string]any `json:"data"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return response.Error(e, err)
	}

	err = h.game.UseItem(e.Request.Context(), e.App, e.Auth.Id, data.ItemId, data.Data)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, nil)
}

func (h *Handlers) DropItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"item_id"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return response.Error(e, err)
	}

	err = h.game.DropItem(e.Request.Context(), e.App, e.Auth.Id, data.ItemId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, nil)
}

func (h *Handlers) GetAvailableActions(e *core.RequestEvent) error {
	res, err := h.game.GetAvailableActions(e.Request.Context(), e.Auth.Id)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) GetEffectView(e *core.RequestEvent) error {
	req := struct {
		EffectId string `json:"effect_id"`
	}{}

	if err := e.BindBody(&req); err != nil {
		return response.Error(e, err)
	}

	res, err := h.game.GetEffectView(e.Request.Context(), e.Auth.Id, req.EffectId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) GetActionView(e *core.RequestEvent) error {
	action := e.Request.URL.Query().Get("action")

	res, err := h.game.GetActionView(e.Request.Context(), e.Auth.Id, model.ActionType(action))
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) RefreshShopHandler(e *core.RequestEvent) error {
	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRefreshShop, nil)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) CurrentSeason(e *core.RequestEvent) error {
	res, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}

func (h *Handlers) IsEventEnded(e *core.RequestEvent) error {
	res, err := h.settings.IsEventEnded(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}
