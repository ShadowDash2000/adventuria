package http

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/actions/custom/buy"
	"adventuria/internal/adventuria_new/actions/custom/done"
	"adventuria/internal/adventuria_new/actions/custom/drop"
	"adventuria/internal/adventuria_new/actions/custom/reroll"
	"adventuria/internal/adventuria_new/actions/custom/update_review"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/event_stats"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
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
	EventStats(ctx context.Context) (*event_stats.EventStatsData, error)
	IsActionsBlocked(ctx context.Context) (bool, error)
}

type result struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type Handlers struct {
	Game game
}

func New(game game) *Handlers {
	return &Handlers{Game: game}
}

func getLang(e *core.RequestEvent) string {
	lang := e.Request.Header.Get("Accept-Language")
	if lang == "" {
		lang = "ru"
	}
	return lang
}

func RespondWithError(e *core.RequestEvent, err error) error {
	lang := getLang(e)

	if appErr, ok := errors.AsType[*errs.AppError](err); ok {
		res := result{
			Success: false,
			Data:    nil,
			Error:   appErr.Code,
			Message: appErr.Message,
		}

		status := http.StatusInternalServerError
		if appErr.Status > 0 {
			status = appErr.Status
		}

		if msg, ok := appErr.Translates[lang]; ok {
			res.Message = msg
		}

		return e.JSON(status, res)
	}

	return e.JSON(http.StatusInternalServerError, result{
		Success: false,
		Error:   "internal_server_error",
	})
}

func RespondWithSuccess(e *core.RequestEvent, data any) error {
	return e.JSON(http.StatusOK, result{
		Success: true,
		Data:    data,
	})
}

func (h *Handlers) UpdateReviewHandler(e *core.RequestEvent) error {
	req := update_review.Request{}

	if err := e.BindBody(&req); err != nil {
		return RespondWithError(e, err)
	}

	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeUpdateReview, req)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollDice, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	req := reroll.Request{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeReroll, req)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	req := drop.Request{}

	err := e.BindBody(&req)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeDrop, req)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	req := done.Request{}

	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeDone, req)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RollWheelHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollWheel, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RollItemHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollItem, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RollItemOnCellHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRollItemOnCell, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) BuyItemHandler(e *core.RequestEvent) error {
	req := buy.Request{}

	err := e.BindBody(&req)
	if err != nil {
		return RespondWithError(e, err)
	}

	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeBuy, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) UseItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string         `json:"item_id"`
		Data   map[string]any `json:"data"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return RespondWithError(e, err)
	}

	err = h.Game.UseItem(e.Request.Context(), e.App, e.Auth.Id, data.ItemId, data.Data)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, nil)
}

func (h *Handlers) DropItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"item_id"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		return RespondWithError(e, err)
	}

	err = h.Game.DropItem(e.Request.Context(), e.App, e.Auth.Id, data.ItemId)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, nil)
}

func (h *Handlers) GetAvailableActions(e *core.RequestEvent) error {
	res, err := h.Game.GetAvailableActions(e.Request.Context(), e.Auth.Id)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) GetEffectView(e *core.RequestEvent) error {
	req := struct {
		EffectId string `json:"effect_id"`
	}{}

	if err := e.BindBody(&req); err != nil {
		return RespondWithError(e, err)
	}

	res, err := h.Game.GetEffectView(e.Request.Context(), e.Auth.Id, req.EffectId)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) GetActionView(e *core.RequestEvent) error {
	action := e.Request.URL.Query().Get("action")

	res, err := h.Game.GetActionView(e.Request.Context(), e.Auth.Id, model.ActionType(action))
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) RefreshShopHandler(e *core.RequestEvent) error {
	res, err := h.Game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeRefreshShop, nil)
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}

func (h *Handlers) EventStats(e *core.RequestEvent) error {
	res, err := h.Game.EventStats(e.Request.Context())
	if err != nil {
		return RespondWithError(e, err)
	}

	return RespondWithSuccess(e, res)
}
