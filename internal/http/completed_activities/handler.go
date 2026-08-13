package completed_activities

import (
	"adventuria/internal/adventuria/completed_activities"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type completedActivities interface {
	GetCompletedActivitiesByCellID(ctx context.Context, cellId string) ([]*completed_activities.CompletedActivity, error)
}

type Handler struct {
	completedActivities completedActivities
}

func NewHandler(completedActivities completedActivities) *Handler {
	return &Handler{completedActivities: completedActivities}
}

func (h *Handler) GetCompletedActivitiesByCellID(e *core.RequestEvent) error {
	cellId := e.Request.URL.Query().Get("cellId")

	res, err := h.completedActivities.GetCompletedActivitiesByCellID(e.Request.Context(), cellId)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, completedActivitiesToView(res))
}
