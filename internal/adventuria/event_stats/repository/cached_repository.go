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
	mu          sync.RWMutex
	repository  dbRepository
	cachedStats map[string]*statCacheEntry
}

type statCacheEntry struct {
	stats     *event_stats.EventStatsEntries
	expiresAt time.Time
}

func NewCachedRepository(repository dbRepository) *CachedRepository {
	return &CachedRepository{
		repository:  repository,
		cachedStats: make(map[string]*statCacheEntry),
	}
}

func (c *CachedRepository) ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error) {
	if cachedStats, ok := c.getCachedStats(seasonId); ok {
		return &event_stats.EventStatsData{
			NextUpdateAt: cachedStats.expiresAt,
			Stats:        cachedStats.stats.Clone(),
		}, nil
	}

	stats, err := c.repository.ComputeStats(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	cachedStats := &statCacheEntry{
		stats:     stats.Clone(),
		expiresAt: time.Now().Add(time.Hour * 3),
	}
	c.cachedStats[seasonId] = cachedStats
	c.mu.Unlock()

	return &event_stats.EventStatsData{
		NextUpdateAt: cachedStats.expiresAt,
		Stats:        stats,
	}, nil
}

func (c *CachedRepository) getCachedStats(seasonId string) (*statCacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cachedStats, ok := c.cachedStats[seasonId]
	if !ok {
		return nil, false
	}

	if cachedStats.expiresAt.Before(time.Now()) {
		delete(c.cachedStats, seasonId)
		return nil, false
	}

	return cachedStats, true
}
