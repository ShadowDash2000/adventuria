package inventory

import (
	inventoryModel "adventuria/internal/adventuria/inventories/model"
	"adventuria/internal/http/response"
	"context"
	"errors"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
}

type inventories interface {
	GetPlayerInventory(ctx context.Context, playerId, seasonId string) (*inventoryModel.PlayerInventory, error)
}

type Handler struct {
	settings    settings
	inventories inventories
}

func NewHandler(settings settings, inventories inventories) *Handler {
	return &Handler{
		settings:    settings,
		inventories: inventories,
	}
}

func (h *Handler) GetInventoryByPlayerID(e *core.RequestEvent) error {
	playerId := e.Request.URL.Query().Get("playerId")

	if playerId == "" {
		return response.Error(e, errors.New("playerId is required"))
	}

	currentSeason, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	res, err := h.inventories.GetPlayerInventory(e.Request.Context(), playerId, currentSeason)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, playerInventoryToView(res))
}
