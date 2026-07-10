package event_stats

import "context"

type repository interface {
	ComputeStats(ctx context.Context, seasonId string) (*EventStatsData, error)
}

type EventStats struct {
	repository repository
}

func NewEventStats(repository repository) *EventStats {
	return &EventStats{repository: repository}
}

func (e *EventStats) ComputeStats(ctx context.Context, seasonId string) (*EventStatsData, error) {
	return e.repository.ComputeStats(ctx, seasonId)
}
