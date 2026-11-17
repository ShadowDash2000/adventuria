package cells

import (
	cellsView "adventuria/internal/adventuria/cells/service/view"
	worldsView "adventuria/internal/adventuria/worlds/service/view"
	"adventuria/internal/http/response"
	"context"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type cellsService interface {
	GetAllBoard(ctx context.Context) ([]*cellsView.CellBoard, error)
	GetDetailedByID(ctx context.Context, id string) (*cellsView.CellDetailed, error)
}

type worldsService interface {
	GetAllBoard(ctx context.Context) ([]*worldsView.WorldBoard, error)
}

type Handler struct {
	cells  cellsService
	worlds worldsService
}

func NewHandler(cells cellsService, worlds worldsService) *Handler {
	return &Handler{
		cells:  cells,
		worlds: worlds,
	}
}

func (h *Handler) GetAllBoard(e *core.RequestEvent) error {
	cells, err := h.cells.GetAllBoard(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	worlds, err := h.worlds.GetAllBoard(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, struct {
		Cells  []*cellBoardView  `json:"cells"`
		Worlds []*worldBoardView `json:"worlds"`
	}{
		Cells:  cellsBoardToView(cells),
		Worlds: worldsBoardToView(worlds),
	})
}

func (h *Handler) GetDetailedByID(e *core.RequestEvent) error {
	cellId := e.Request.URL.Query().Get("cellId")

	cell, err := h.cells.GetDetailedByID(e.Request.Context(), cellId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, cellDetailedToView(cell))
}
