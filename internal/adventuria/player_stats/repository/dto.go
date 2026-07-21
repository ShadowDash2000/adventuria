package repository

import "adventuria/internal/adventuria/model"

type activitiesStatsDTO struct {
	GamesCompleted   int `json:"games_completed"`
	MoviesCompleted  int `json:"movies_completed"`
	GymsCompleted    int `json:"gyms_completed"`
	KaraokeCompleted int `json:"karaoke_completed"`
}

func activitiesStatsToDTO(stats model.ActivitiesStats) activitiesStatsDTO {
	return activitiesStatsDTO{
		GamesCompleted:   stats.GamesCompleted,
		MoviesCompleted:  stats.MoviesCompleted,
		GymsCompleted:    stats.GymsCompleted,
		KaraokeCompleted: stats.KaraokeCompleted,
	}
}

func activitiesStatsFromDTO(dto activitiesStatsDTO) model.ActivitiesStats {
	return model.ActivitiesStats{
		GamesCompleted:   dto.GamesCompleted,
		MoviesCompleted:  dto.MoviesCompleted,
		GymsCompleted:    dto.GymsCompleted,
		KaraokeCompleted: dto.KaraokeCompleted,
	}
}
