package model

import (
	"errors"
	"time"
)

type SeasonData struct {
	Id              string
	Name            string
	Slug            string
	SeasonDateStart time.Time
	SeasonDateEnd   time.Time
}

type Season struct {
	data  SeasonData
	isNew bool
}

type SeasonCreate struct {
	Name            string
	Slug            string
	SeasonDateStart time.Time
	SeasonDateEnd   time.Time
}

func NewSeason(data SeasonCreate) (*Season, error) {
	if data.Name == "" {
		return nil, errors.New("season: name is empty")
	}
	if data.Slug == "" {
		return nil, errors.New("season: slug is empty")
	}
	if data.SeasonDateStart.IsZero() {
		return nil, errors.New("season: season date start is empty")
	}
	if data.SeasonDateEnd.IsZero() {
		return nil, errors.New("season: season date end is empty")
	}

	return &Season{
		data: SeasonData{
			Name:            data.Name,
			Slug:            data.Slug,
			SeasonDateStart: data.SeasonDateStart,
			SeasonDateEnd:   data.SeasonDateEnd,
		},
		isNew: true,
	}, nil
}

func RestoreSeason(data SeasonData) *Season {
	return &Season{
		data:  data,
		isNew: false,
	}
}

func (s *Season) IsNew() bool {
	return s.isNew
}

func (s *Season) ID() string {
	return s.data.Id
}

func (s *Season) Name() string {
	return s.data.Name
}

func (s *Season) Slug() string {
	return s.data.Slug
}

func (s *Season) SeasonDateStart() time.Time {
	return s.data.SeasonDateStart
}

func (s *Season) SeasonDateEnd() time.Time {
	return s.data.SeasonDateEnd
}
