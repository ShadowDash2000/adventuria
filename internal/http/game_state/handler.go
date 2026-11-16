package game_state

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/npcs"
	"adventuria/internal/adventuria/player_progress"
	"adventuria/internal/http/response"
	"context"
	"errors"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
	GetNextDayBySeasonID(ctx context.Context, seasonId string) (time.Time, error)
}

type playerInfo interface {
	GetByID(ctx context.Context, id string) (*model.PlayerInfo, error)
}

type playerProgress interface {
	GetFirstOrDefault(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error)
	GetActivePlayersBySeasonID(ctx context.Context, seasonId string) ([]*player_progress.ActivePlayer, error)
}

type actorsService interface {
	GetActiveActorsBySeasonID(ctx context.Context, seasonId string) ([]*actors.ActiveActor, error)
}

type npcsService interface {
	GetActiveNPCs(ctx context.Context) ([]*npcs.ActiveNPC, error)
}

type Handler struct {
	settings       settings
	playerInfo     playerInfo
	playerProgress playerProgress
	actors         actorsService
	npc            npcsService
}

func NewHandler(
	settings settings,
	playerInfo playerInfo,
	playerProgress playerProgress,
	actors actorsService,
	npc npcsService,
) *Handler {
	return &Handler{
		settings:       settings,
		playerInfo:     playerInfo,
		playerProgress: playerProgress,
		actors:         actors,
		npc:            npc,
	}
}

func (h *Handler) GetGameState(e *core.RequestEvent) error {
	currentSeason, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	nextDayAt, err := h.settings.GetNextDayBySeasonID(e.Request.Context(), currentSeason)
	if err != nil && !errors.Is(err, errs.ErrSeasonEnded) {
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

	return response.Success(e, gameStateToView(currentSeason, nextDayAt, player, progress))
}

func (h *Handler) GetActiveActorsList(e *core.RequestEvent) error {
	currentSeason, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	activePlayers, err := h.playerProgress.GetActivePlayersBySeasonID(e.Request.Context(), currentSeason)
	if err != nil {
		return response.Error(e, err)
	}

	activeNPCs, err := h.npc.GetActiveNPCs(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, struct {
		Players []*activePlayerView `json:"players"`
		NPCs    []*activeNPCView    `json:"npcs"`
	}{
		Players: activePlayersToView(activePlayers),
		NPCs:    activeNPCsToView(activeNPCs),
	})
}

func (h *Handler) GetActiveActors(e *core.RequestEvent) error {
	currentSeason, err := h.settings.CurrentSeason(e.Request.Context())
	if err != nil {
		return response.Error(e, err)
	}

	activeActors, err := h.actors.GetActiveActorsBySeasonID(e.Request.Context(), currentSeason)
	if err != nil {
		return response.Error(e, err)
	}

	return response.Success(e, activeActorsToView(activeActors))
}
