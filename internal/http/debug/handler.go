package debug

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/actions/custom/move_to_cell_id"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type game interface {
	DoAction(ctx context.Context, pb core.App, playerId string, actionType model.ActionType, req model.ActionRequest) (any, error)
}

type Handler struct {
	game game
}

func NewHandler(game game) *Handler {
	return &Handler{
		game: game,
	}
}
func (h *Handler) MoveToCellID(e *core.RequestEvent) error {
	req := move_to_cell_id.Request{}

	err := e.BindBody(&req)
	if err != nil {
		return response.Error(e, err)
	}

	res, err := h.game.DoAction(e.Request.Context(), e.App, e.Auth.Id, actions.ActionTypeMoveToCellId, req)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, res)
}
