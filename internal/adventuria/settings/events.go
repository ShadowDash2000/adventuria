package settings

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
)

type OnNewDayEvent struct {
	event.Event
	Settings   *model.Settings
	Season     *model.Season
	CurrentDay int
	DaysPassed int
}

func (s *Settings) OnNewDay() *event.Hook[*OnNewDayEvent] {
	return s.onNewDay
}
