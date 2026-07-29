package repository

import (
	"adventuria/internal/adventuria/event_stats"
	"context"
	"sync"
	"time"
)

type dbRepository interface {
	ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsEntries, error)
}

type CachedRepository struct {
	mu              sync.RWMutex
	repository      dbRepository
	cachedStats     *event_stats.EventStatsEntries
	cachedExpiresAt time.Time
}

func NewCachedRepository(repository dbRepository) *CachedRepository {
	return &CachedRepository{
		repository: repository,
	}
}

func (c *CachedRepository) ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error) {
	if cachedStats, ok := c.getCachedStats(); ok {
		return &event_stats.EventStatsData{
			NextUpdateAt: c.cachedExpiresAt,
			Stats:        cachedStats.Clone(),
		}, nil
	}

	stats, err := c.repository.ComputeStats(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cachedStats = stats.Clone()
	c.cachedExpiresAt = time.Now().Add(time.Hour * 3)
	c.mu.Unlock()

	return &event_stats.EventStatsData{
		NextUpdateAt: c.cachedExpiresAt,
		Stats:        stats,
	}, nil
}

func (c *CachedRepository) getCachedStats() (*event_stats.EventStatsEntries, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cachedStats == nil {
		return nil, false
	}

	if c.cachedExpiresAt.Before(time.Now()) {
		c.cachedStats = nil
		return nil, false
	}

	return c.cachedStats, true
}
