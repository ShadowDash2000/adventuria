package seasons

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, season *model.Season) (*model.Season, error)
	GetFirst(ctx context.Context) (*model.Season, error)
	GetByID(ctx context.Context, id string) (*model.Season, error)
}

type Seasons struct {
	repository repository
}

var defaultSeason = model.SeasonCreate{
	Name:            "Season 1",
	Slug:            "season-1",
	SeasonDateStart: time.Now(),
	SeasonDateEnd:   time.Now().AddDate(0, 1, 0),
}

func NewSeasons(repo repository) *Seasons {
	return &Seasons{repository: repo}
}

func (s *Seasons) GetFirstOrDefault(ctx context.Context) (*model.Season, error) {
	season, err := s.repository.GetFirst(ctx)
	if err == nil {
		return season, err
	} else if !errors.Is(err, errs.ErrSeasonNotFound) {
		return nil, err
	}

	season, err = model.NewSeason(uuid.New(), defaultSeason)
	if err != nil {
		return nil, err
	}

	season, err = s.repository.Create(ctx, season)
	if err != nil {
		return nil, err
	}

	return season, nil
}

func (s *Seasons) GetByID(ctx context.Context, id string) (*model.Season, error) {
	return s.repository.GetByID(ctx, id)
}
