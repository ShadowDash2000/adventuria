package repository

import (
	"adventuria/internal/adventuria/event_stats"
	"adventuria/pkg/cache"
	"context"
	"time"
)

type dbRepository interface {
	ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error)
}

type CachedRepository struct {
	repository dbRepository
	cache      cache.Cache[string, any]
}

func NewCachedRepository(repository dbRepository) *CachedRepository {
	return &CachedRepository{
		repository: repository,
		cache:      cache.NewMemoryCache[string, any](time.Hour*24, false),
	}
}

// ComputeStats TODO cache
func (c *CachedRepository) ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error) {
	return c.repository.ComputeStats(ctx, seasonId)
}
