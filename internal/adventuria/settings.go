package adventuria

import (
	"adventuria/pkg/event"
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Settings struct {
	core.BaseRecordProxy

	onKillParser *event.Hook[*OnKillParserEvent]
}

type OnKillParserEvent struct {
	event.Event
}

func NewSettings() (*Settings, error) {
	s := &Settings{}

	if err := s.init(); err != nil {
		return nil, err
	}
	s.initHooks()
	s.bindHooks()
	s.RegisterSettingsCron()

	return s, nil
}

func DefaultSettings() (*core.Record, error) {
	record := core.NewRecord(GameCollections.Get(CollectionSettings))
	record.Set("eventDateStart", types.NowDateTime())
	record.Set("currentWeek", 0)
	record.Set("timerTimeLimit", 14400)
	record.Set("limitExceedPenalty", 2)
	record.Set("pointsForDrop", -2)
	record.Set("dropsToJail", 2)
	record.Set("igdbParseSettings", "game_type = 0 & platforms = {6}")
	return record, nil
}

func (s *Settings) initHooks() {
	s.onKillParser = &event.Hook[*OnKillParserEvent]{}
}

func (s *Settings) bindHooks() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	PocketBase.OnRecordUpdate(CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		if ok := e.Record.GetBool("kill_parser"); ok {
			_ = s.onKillParser.Trigger(&OnKillParserEvent{})
			e.Record.Set("kill_parser", false)
		}
		return e.Next()
	})
}

func (s *Settings) init() error {
	record, err := PocketBase.FindFirstRecordByFilter(
		CollectionSettings,
		"",
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if record != nil {
		s.SetProxyRecord(record)
	} else {
		record, err = DefaultSettings()
		if err != nil {
			return err
		}

		s.SetProxyRecord(record)
		err = PocketBase.Save(s)
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
	PocketBase.Cron().MustAdd("settings", "* * * * *", func() {
		week := s.GetCurrentWeekNum()
		if s.CurrentWeek() == week {
			return
		}

		s.SetCurrentWeek(week)
		err := PocketBase.Save(s)
		if err != nil {
			PocketBase.Logger().Error("save settings failed", "err", err)
		}

		err = ResetAllTimers(s.TimerTimeLimit(), s.LimitExceedPenalty())
		if err != nil {
			PocketBase.Logger().Error("failed to clear timers", "err", err)
		}
	})
}

func (s *Settings) IGDBGamesParsed() uint64 {
	return uint64(s.GetInt("igdb_games_parsed"))
}

func (s *Settings) SetIGDBGamesParsed(count uint64) {
	s.Set("igdb_games_parsed", count)
}

func (s *Settings) DisableIGDBParser() bool {
	return s.GetBool("disable_igdb_parser")
}

func (s *Settings) DisableSteamParser() bool {
	return s.GetBool("disable_steam_parser")
}

func (s *Settings) DisableCheapsharkParser() bool {
	return s.GetBool("disable_cheapshark_parser")
}

func (s *Settings) DisableHLTBParser() bool {
	return s.GetBool("disable_hltb_parser")
}

func (s *Settings) OnKillParser() *event.Hook[*OnKillParserEvent] {
	return s.onKillParser
}
