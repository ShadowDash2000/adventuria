package game_state

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/http/response"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
}

type playerInfo interface {
	GetByID(ctx context.Context, id string) (*model.PlayerInfo, error)
}

type playerProgress interface {
	GetFirstOrDefault(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error)
}

type Handler struct {
	settings       settings
	playerInfo     playerInfo
	playerProgress playerProgress
}

func NewHandler(settings settings, playerInfo playerInfo, playerProgress playerProgress) *Handler {
	return &Handler{
		settings:       settings,
		playerInfo:     playerInfo,
		playerProgress: playerProgress,
	}
}

func (h *Handler) GetGameState(e *core.RequestEvent) error {
	currentSeason, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	var (
		player   *model.PlayerInfo
		progress *model.PlayerProgress
	)
	if e.Auth != nil {
		player, err = h.playerInfo.GetByID(e.Request.Context(), e.Auth.Id)
		if err != nil {
			return response.Error(e, err)
		}

		progress, err = h.playerProgress.GetFirstOrDefault(e.Request.Context(), e.Auth.Id, currentSeason)
		if err != nil {
			return response.Error(e, err)
		}
	}

	return response.Success(e, stateToView(currentSeason, player, progress))
}
