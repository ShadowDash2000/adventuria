package settings

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, settings *model.Settings) (*model.Settings, error)
	GetFirst(ctx context.Context) (*model.Settings, error)
	IsActionsBlocked(ctx context.Context) (bool, error)
	UpdateIGDBGamesParsedByID(ctx context.Context, id string, amount int) error
}

type seasons interface {
	GetFirstOrDefault(ctx context.Context) (*model.Season, error)
}

type Settings struct {
	repository repository
	seasons    seasons
}

func defaultSettings(season string) model.SettingsCreate {
	return model.SettingsCreate{
		CurrentSeason: season,
		DropsToJail:   2,
	}
}

func NewSettings(repo repository, seasons seasons) *Settings {
	return &Settings{
		repository: repo,
		seasons:    seasons,
	}
}

func (s *Settings) GetFirstOrDefault(ctx context.Context) (*model.Settings, error) {
	settings, err := s.repository.GetFirst(ctx)
	if err == nil {
		return settings, err
	} else if !errors.Is(err, errs.ErrSettingsNotFound) {
		return nil, err
	}

	season, err := s.seasons.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	settings, err = model.NewSettings(uuid.New(), defaultSettings(season.ID()))
	if err != nil {
		return nil, err
	}

	settings, err = s.repository.Create(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *Settings) IsActionsBlocked(ctx context.Context) (bool, error) {
	return s.repository.IsActionsBlocked(ctx)
}

func (s *Settings) UpdateIGDBGamesParsedByID(ctx context.Context, id string, amount int) error {
	return s.repository.UpdateIGDBGamesParsedByID(ctx, id, amount)
}
