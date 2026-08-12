package completed_activities

import (
	"adventuria/internal/adventuria/players"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type game interface {
	GetCompletedActivitiesByCellID(ctx context.Context, cellId string) ([]*players.CompletedActivity, error)
}

type Handler struct {
	game game
}

func New(game game) *Handler {
	return &Handler{game: game}
}

func (h *Handler) GetCompletedActivitiesByCellID(e *core.RequestEvent) error {
	cellId := e.Request.URL.Query().Get("cellId")

	res, err := h.game.GetCompletedActivitiesByCellID(e.Request.Context(), cellId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, completedActivitiesToView(res))
}
