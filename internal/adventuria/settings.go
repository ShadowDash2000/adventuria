package adventuria

import (
	"adventuria/pkg/collections"
	"database/sql"
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Settings struct {
	core.BaseRecordProxy
	app core.App
}

func NewSettings(cols *collections.Collections, app core.App) (*Settings, error) {
	s := &Settings{app: app}

	record, err := s.app.FindFirstRecordByFilter(
		TableSettings,
		"",
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if record != nil {
		s.SetProxyRecord(record)
	} else {
		record, err = DefaultSettings(cols)
		if err != nil {
			return nil, err
		}

		s.SetProxyRecord(record)
		err = app.Save(s)
		if err != nil {
			return nil, err
		}
	}

	s.bindHooks()
	s.RegisterSettingsCron()

	return s, nil
}

func DefaultSettings(cols *collections.Collections) (*core.Record, error) {
	collection, err := cols.Get(TableSettings)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Set("eventDateStart", types.NowDateTime())
	record.Set("currentWeek", 0)
	record.Set("timerTimeLimit", 14400)
	return record, nil
}

func (s *Settings) bindHooks() {
	s.app.OnRecordAfterCreateSuccess(TableSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	s.app.OnRecordAfterUpdateSuccess(TableSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
}

func (s *Settings) EventDateStart() types.DateTime {
	return s.GetDateTime("eventDateStart")
}

func (s *Settings) CurrentWeek() int {
	return s.GetInt("currentWeek")
}

func (s *Settings) SetCurrentWeek(w int) {
	s.Set("currentWeek", w)
}

func (s *Settings) getCurrentWeekNum() int {
	daysPassed := int(types.NowDateTime().Sub(s.EventDateStart()).Hours() / 24)
	return (daysPassed / 7) + 1
}

func (s *Settings) TimerTimeLimit() int {
	return s.GetInt("timerTimeLimit")
}

func (s *Settings) BlockAllActions() bool {
	return s.GetBool("blockAllActions")
}

func (s *Settings) Rules() {

}

func (s *Settings) RegisterSettingsCron() {
	s.app.Cron().MustAdd("settings", "* * * * *", func() {
		week := s.getCurrentWeekNum()
		if s.CurrentWeek() == week {
			return
		}

		s.SetCurrentWeek(week)
		err := s.app.Save(s)
		if err != nil {
			s.app.Logger().Error("save settings failed", "err", err)
		}

		err = ResetAllTimers(s.TimerTimeLimit(), s.app)
		if err != nil {
			s.app.Logger().Error("failed to clear timers", "err", err)
		}
	})
}
