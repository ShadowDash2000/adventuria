package update_review

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type reviewsService interface {
	UpdateByActionAndPlayerID(ctx context.Context, actionId, playerId string, input reviews.UpdateInput) (*model.Review, error)
}

type Handler struct {
	reviews reviewsService
}

func NewHandler(reviews reviewsService) *Handler {
	return &Handler{reviews: reviews}
}

type request struct {
	ActionID string   `json:"action_id" form:"action_id"`
	Comment  *string  `json:"comment" form:"comment"`
	Score    *float64 `json:"score" form:"score"`
}

func (h *Handler) UpdateReviewByActionID(e *core.RequestEvent) error {
	req := request{}
	err := e.BindBody(&req)
	if err != nil {
		return response.Error(e, err)
	}

	res, err := h.reviews.UpdateByActionAndPlayerID(e.Request.Context(), req.ActionID, e.Auth.Id, reviews.UpdateInput{
		Comment: req.Comment,
		Score:   req.Score,
	})
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}
