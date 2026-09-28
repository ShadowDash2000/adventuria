package debug

import (
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type game interface {
	MoveToCellID(ctx context.Context, pb core.App, playerId, cellId string) error
	AddItemByID(ctx context.Context, pb core.App, playerId, itemId string) error
}

type Handler struct {
	game game
}

func NewHandler(game game) *Handler {
	return &Handler{
		game: game,
	}
}

type moveToCellIDRequest struct {
	PlayerId string `json:"player_id"`
	CellId   string `json:"cell_id"`
}

func (h *Handler) MoveToCellID(e *core.RequestEvent) error {
	req := moveToCellIDRequest{}

	err := e.BindBody(&req)
	if err != nil {
		return response.Error(e, err)
	}

	err = h.game.MoveToCellID(e.Request.Context(), e.App, req.PlayerId, req.CellId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, nil)
}

type addItemByIDRequest struct {
	PlayerId string `json:"player_id"`
	ItemId   string `json:"item_id"`
}

func (h *Handler) AddItemByID(e *core.RequestEvent) error {
	req := addItemByIDRequest{}

	err := e.BindBody(&req)
	if err != nil {
		return response.Error(e, err)
	}

	err = h.game.AddItemByID(e.Request.Context(), e.App, req.PlayerId, req.ItemId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, nil)
}
