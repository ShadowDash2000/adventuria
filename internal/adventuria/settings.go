package adventuria

import (
	"adventuria/pkg/collections"
	"database/sql"
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Settings struct {
	core.BaseRecordProxy
	app  core.App
	cols *collections.Collections
}

func NewSettings(cols *collections.Collections, app core.App) *Settings {
	s := &Settings{
		app:  app,
		cols: cols,
	}

	s.init()
	s.bindHooks()
	s.RegisterSettingsCron()

	return s
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
	record.Set("limitExceedPenalty", 2)
	record.Set("pointsForDrop", -2)
	record.Set("dropsToJail", 2)
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

func (s *Settings) init() error {
	record, err := s.app.FindFirstRecordByFilter(
		TableSettings,
		"",
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if record != nil {
		s.SetProxyRecord(record)
	} else {
		record, err = DefaultSettings(s.cols)
		if err != nil {
			return err
		}

		s.SetProxyRecord(record)
		err = s.app.Save(s)
		if err != nil {
			return err
		}
	}

	return nil
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

func (s *Settings) DaysPassedFromEventStart() int {
	return int(types.NowDateTime().Sub(s.EventDateStart()).Hours() / 24)
}

func (s *Settings) GetCurrentWeekNum() int {
	return (s.DaysPassedFromEventStart() / 7) + 1
}

func (s *Settings) TimerTimeLimit() int {
	return s.GetInt("timerTimeLimit")
}

func (s *Settings) LimitExceedPenalty() int {
	return s.GetInt("limitExceedPenalty")
}

func (s *Settings) BlockAllActions() bool {
	return s.GetBool("blockAllActions")
}

func (s *Settings) PointsForDrop() int {
	return s.GetInt("pointsForDrop")
}

func (s *Settings) DropsToJail() int {
	return s.GetInt("dropsToJail")
}

func (s *Settings) NextTimerResetDate() types.DateTime {
	weeks := time.Duration(s.CurrentWeek()*7*24) * time.Hour
	return s.EventDateStart().Add(weeks)
}

func (s *Settings) CheckActionsBlock() *hook.Handler[*core.RequestEvent] {
	return &hook.Handler[*core.RequestEvent]{
		Id:   "settingsCheckActionsBlock",
		Func: s.checkActionsBlock(),
	}
}

func (s *Settings) checkActionsBlock() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if s.BlockAllActions() {
			return e.ForbiddenError("All actions are temporarily blocked", nil)
		}

		return e.Next()
	}
}

func (s *Settings) RegisterSettingsCron() {
	s.app.Cron().MustAdd("settings", "* * * * *", func() {
		week := s.GetCurrentWeekNum()
		if s.CurrentWeek() == week {
			return
		}

		s.SetCurrentWeek(week)
		err := s.app.Save(s)
		if err != nil {
			s.app.Logger().Error("save settings failed", "err", err)
		}

		err = ResetAllTimers(s.TimerTimeLimit(), s.LimitExceedPenalty(), s.app)
		if err != nil {
			s.app.Logger().Error("failed to clear timers", "err", err)
		}
	})
}
