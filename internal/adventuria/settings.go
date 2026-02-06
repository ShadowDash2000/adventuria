package adventuria

import (
	"adventuria/internal/adventuria/schema"
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

func NewSettings(ctx AppContext) (*Settings, error) {
	s := &Settings{}

	if err := s.init(ctx); err != nil {
		return nil, err
	}
	s.initHooks()
	s.bindHooks(ctx)
	s.RegisterSettingsCron(ctx)

	return s, nil
}

func DefaultSettings() (*core.Record, error) {
	record := core.NewRecord(GameCollections.Get(schema.CollectionSettings))
	record.Set(schema.SettingsSchema.EventDateStart, types.NowDateTime())
	record.Set(schema.SettingsSchema.CurrentWeek, 0)
	record.Set(schema.SettingsSchema.TimerTimeLimit, 14400)
	record.Set(schema.SettingsSchema.LimitExceedPenalty, 2)
	record.Set(schema.SettingsSchema.PointsForDrop, -2)
	record.Set(schema.SettingsSchema.DropsToJail, 2)
	return record, nil
}

func (s *Settings) initHooks() {
	s.onKillParser = &event.Hook[*OnKillParserEvent]{}
}

func (s *Settings) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordUpdate(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		if ok := e.Record.GetBool(schema.SettingsSchema.KillParser); ok {
			_, err := s.onKillParser.Trigger(&OnKillParserEvent{})
			if err != nil {
				e.App.Logger().Error("Failed to trigger kill parser event", "err", err)
			}
			e.Record.Set(schema.SettingsSchema.KillParser, false)
		}
		return e.Next()
	})
}

func (s *Settings) init(ctx AppContext) error {
	record, err := ctx.App.FindFirstRecordByFilter(
		schema.CollectionSettings,
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
		err = ctx.App.Save(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Settings) EventDateStart() types.DateTime {
	return s.GetDateTime(schema.SettingsSchema.EventDateStart)
}

func (s *Settings) CurrentWeek() int {
	return s.GetInt(schema.SettingsSchema.CurrentWeek)
}

func (s *Settings) SetCurrentWeek(w int) {
	s.Set(schema.SettingsSchema.CurrentWeek, w)
}

func (s *Settings) DaysPassedFromEventStart() int {
	return int(types.NowDateTime().Sub(s.EventDateStart()).Hours() / 24)
}

func (s *Settings) GetCurrentWeekNum() int {
	return (s.DaysPassedFromEventStart() / 7) + 1
}

func (s *Settings) TimerTimeLimit() int {
	return s.GetInt(schema.SettingsSchema.TimerTimeLimit)
}

func (s *Settings) LimitExceedPenalty() int {
	return s.GetInt(schema.SettingsSchema.LimitExceedPenalty)
}

func (s *Settings) BlockAllActions() bool {
	return s.GetBool(schema.SettingsSchema.BlockAllActions)
}

func (s *Settings) PointsForDrop() int {
	return s.GetInt(schema.SettingsSchema.PointsForDrop)
}

func (s *Settings) DropsToJail() int {
	return s.GetInt(schema.SettingsSchema.DropsToJail)
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

func (s *Settings) RegisterSettingsCron(ctx AppContext) {
	ctx.App.Cron().MustAdd("settings", "* * * * *", func() {
		week := s.GetCurrentWeekNum()
		if s.CurrentWeek() == week {
			return
		}

		s.SetCurrentWeek(week)
		err := ctx.App.Save(s)
		if err != nil {
			ctx.App.Logger().Error("save settings failed", "err", err)
		}

		err = ResetAllTimers(AppContext{App: PocketBase}, s.TimerTimeLimit(), s.LimitExceedPenalty())
		if err != nil {
			ctx.App.Logger().Error("failed to clear timers", "err", err)
		}
	})
}

func (s *Settings) IGDBGamesParsed() uint64 {
	return uint64(s.GetInt(schema.SettingsSchema.IgdbGamesParsed))
}

func (s *Settings) SetIGDBGamesParsed(count uint64) {
	s.Set(schema.SettingsSchema.IgdbGamesParsed, count)
}

func (s *Settings) DisableIGDBParser() bool {
	return s.GetBool(schema.SettingsSchema.DisableIgdbParser)
}

func (s *Settings) DisableSteamParser() bool {
	return s.GetBool(schema.SettingsSchema.DisableSteamParser)
}

func (s *Settings) DisableCheapsharkParser() bool {
	return s.GetBool(schema.SettingsSchema.DisableCheapsharkParser)
}

func (s *Settings) DisableHLTBParser() bool {
	return s.GetBool(schema.SettingsSchema.DisableHltbParser)
}

func (s *Settings) OnKillParser() *event.Hook[*OnKillParserEvent] {
	return s.onKillParser
}

func (s *Settings) IgdbForceUpdateGames() bool {
	return s.GetBool(schema.SettingsSchema.IgdbForceUpdateGames)
}
