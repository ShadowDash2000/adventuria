package event_stats

import (
	"adventuria/internal/adventuria/event_stats"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type eventStats interface {
	ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error)
}

type Handler struct {
	eventStats eventStats
}

func NewHandler(eventStats eventStats) *Handler {
	return &Handler{eventStats: eventStats}
}

func (h *Handler) EventStats(e *core.RequestEvent) error {
	seasonId := e.Request.URL.Query().Get("seasonId")

	res, err := h.eventStats.ComputeStats(e.Request.Context(), seasonId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}
