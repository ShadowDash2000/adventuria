package adventuria

import (
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Settings struct {
	core.BaseRecordProxy
	locator SettingsLocator
}

func NewSettings(locator SettingsLocator) *Settings {
	s := &Settings{
		locator: locator,
	}

	s.init()
	s.bindHooks()
	s.RegisterSettingsCron()

	return s
}

func DefaultSettings(locator PocketBaseLocator) (*core.Record, error) {
	collection, err := locator.Collections().Get(TableSettings)
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
	record.Set("igdbParseSettings", "game_type = 0 & platforms = {6}")
	return record, nil
}

func (s *Settings) bindHooks() {
	s.locator.PocketBase().OnRecordAfterCreateSuccess(TableSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	s.locator.PocketBase().OnRecordAfterUpdateSuccess(TableSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
}

func (s *Settings) init() error {
	record, err := s.locator.PocketBase().FindFirstRecordByFilter(
		TableSettings,
		"",
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if record != nil {
		s.SetProxyRecord(record)
	} else {
		record, err = DefaultSettings(s.locator)
		if err != nil {
			return err
		}

		s.SetProxyRecord(record)
		err = s.locator.PocketBase().Save(s)
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

func (s *Settings) TwitchClientID() string {
	return s.GetString("twitchClientId")
}

func (s *Settings) TwitchClientSecret() string {
	return s.GetString("twitchClientSecret")
}

func (s *Settings) IGDBParseSettings() string {
	return s.GetString("igdbParseSettings")
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
	s.locator.PocketBase().Cron().MustAdd("settings", "* * * * *", func() {
		week := s.GetCurrentWeekNum()
		if s.CurrentWeek() == week {
			return
		}

		s.SetCurrentWeek(week)
		err := s.locator.PocketBase().Save(s)
		if err != nil {
			s.locator.PocketBase().Logger().Error("save settings failed", "err", err)
		}

		err = ResetAllTimers(s.locator, s.TimerTimeLimit(), s.LimitExceedPenalty())
		if err != nil {
			s.locator.PocketBase().Logger().Error("failed to clear timers", "err", err)
		}
	})
}
