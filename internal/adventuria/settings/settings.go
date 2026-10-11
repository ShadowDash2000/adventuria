package settings

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
	"time"
)

type repository interface {
	Create(ctx context.Context, settings *model.Settings) (*model.Settings, error)
	GetFirst(ctx context.Context) (*model.Settings, error)
	IsActionsBlocked(ctx context.Context) (bool, error)
	CurrentSeason(ctx context.Context) (string, error)
	IsEventEnded(ctx context.Context) (bool, error)
	ChangeIGDBGamesParsedByID(ctx context.Context, id string, amount int) error
	EnergyDefault(ctx context.Context) (int, error)
	UpdateCurrentDayByID(ctx context.Context, id string, day int) error
}

type seasons interface {
	InitDefault(ctx context.Context) (*model.Season, error)
	GetByID(ctx context.Context, id string) (*model.Season, error)
}

type Settings struct {
	repository repository
	seasons    seasons

	onNewDay *event.Hook[*OnNewDayEvent]
}

func defaultSettings(season string) (*model.Settings, error) {
	settings, err := model.NewSettings(model.SettingsCreate{
		CurrentSeason: season,
		DropsToJail:   2,
	})
	if err != nil {
		return nil, err
	}

	settings.SetDisableIgdbGamesParser(true)

	return settings, nil
}

func NewSettings(repo repository, seasons seasons) *Settings {
	return &Settings{
		repository: repo,
		seasons:    seasons,

		onNewDay: &event.Hook[*OnNewDayEvent]{},
	}
}

func (s *Settings) InitDefault(ctx context.Context) (*model.Settings, error) {
	settings, err := s.repository.GetFirst(ctx)
	if err == nil {
		return settings, err
	} else if !errors.Is(err, errs.ErrSettingsNotFound) {
		return nil, err
	}

	season, err := s.seasons.InitDefault(ctx)
	if err != nil {
		return nil, err
	}

	settings, err = defaultSettings(season.ID())
	if err != nil {
		return nil, err
	}

	settings, err = s.repository.Create(ctx, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *Settings) UpdateCurrentDay(ctx context.Context) error {
	settings, err := s.GetFirst(ctx)
	if err != nil {
		return err
	}

	season, err := s.seasons.GetByID(ctx, settings.CurrentSeason())
	if err != nil {
		return err
	}

	if season.IsEnded(time.Now()) {
		return nil
	}

	newCurrentDay := max(
		int(time.Now().Sub(season.SeasonDateStart()).Hours()/24),
		0,
	)

	if newCurrentDay == settings.CurrentDay() {
		return nil
	}

	err = s.repository.UpdateCurrentDayByID(ctx, settings.ID(), newCurrentDay)
	if err != nil {
		return err
	}

	if newCurrentDay > settings.CurrentDay() {
		err = s.OnNewDay().Trigger(ctx, &OnNewDayEvent{
			Settings:   settings,
			Season:     season,
			CurrentDay: newCurrentDay,
			DaysPassed: newCurrentDay - settings.CurrentDay(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Settings) GetFirst(ctx context.Context) (*model.Settings, error) {
	return s.repository.GetFirst(ctx)
}

func (s *Settings) CurrentSeason(ctx context.Context) (string, error) {
	return s.repository.CurrentSeason(ctx)
}

func (s *Settings) IsActionsBlocked(ctx context.Context) (bool, error) {
	return s.repository.IsActionsBlocked(ctx)
}

func (s *Settings) IsEventEnded(ctx context.Context) (bool, error) {
	return s.repository.IsEventEnded(ctx)
}

func (s *Settings) ChangeIGDBGamesParsedByID(ctx context.Context, id string, amount int) error {
	return s.repository.ChangeIGDBGamesParsedByID(ctx, id, amount)
}

func (s *Settings) EnergyDefault(ctx context.Context) (int, error) {
	return s.repository.EnergyDefault(ctx)
}
