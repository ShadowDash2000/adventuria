package player_stats

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Save(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error)
	GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerStats, error)
}

type PlayerStats struct {
	repository repository
}

func NewPlayerStats(repository repository) *PlayerStats {
	return &PlayerStats{
		repository: repository,
	}
}

func (p *PlayerStats) Save(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error) {
	return p.repository.Save(ctx, stats)
}

func (p *PlayerStats) GetOrCreate(ctx context.Context, playerId, seasonId string) (*model.PlayerStats, error) {
	stats, err := p.repository.GetByPlayerId(ctx, playerId, seasonId)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerStatsNotFound) {
			return model.NewPlayerStats(model.PlayerStatsCreate{
				Player: playerId,
				Season: seasonId,
			})
		}
		return nil, err
	}

	return stats, nil
}
