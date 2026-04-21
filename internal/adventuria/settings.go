package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"database/sql"
	"errors"
	"fmt"

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

	return s, nil
}

func (s *Settings) initHooks() {
	s.onKillParser = &event.Hook[*OnKillParserEvent]{}
}

func (s *Settings) bindHooks(ctx AppContext) {
	ctx.App.OnRecordCreate(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		err := expandSettingsWithSeason(AppContext{App: e.App}, e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		_ = GamePlayers.RefetchAllInMemory(AppContext{App: e.App})
		s.SetProxyRecord(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordUpdate(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		err := expandSettingsWithSeason(AppContext{App: e.App}, e.Record)
		if err != nil {
			return err
		}
		if ok := e.Record.GetBool(schema.SettingsSchema.KillParser); ok {
			_, err = s.onKillParser.Trigger(&OnKillParserEvent{})
			if err != nil {
				e.App.Logger().Error("Failed to trigger kill parser event", "err", err)
			}
			e.Record.Set(schema.SettingsSchema.KillParser, false)
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		if s.CurrentSeason() != e.Record.GetString(schema.SettingsSchema.CurrentSeason) {
			_ = GamePlayers.RefetchAllInMemory(AppContext{App: e.App})
		}
		s.SetProxyRecord(e.Record)
		return e.Next()
	})

	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionSeasons).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == s.CurrentSeason() {
			expand := s.Expand()
			oldSeasonDateStart := expand[schema.SettingsSchema.CurrentSeason].(*core.Record).GetDateTime(schema.SeasonSchema.SeasonDateStart)
			expand[schema.SettingsSchema.CurrentSeason] = e.Record
			s.SetExpand(expand)

			if oldSeasonDateStart != s.CurrentSeasonDateStart() {
				_ = GamePlayers.RefetchAllInMemory(AppContext{App: e.App})
			}
		}
		return e.Next()
	})
}

func expandSettingsWithSeason(ctx AppContext, record *core.Record) error {
	errs := ctx.App.ExpandRecord(record, []string{schema.SettingsSchema.CurrentSeason}, nil)
	if len(errs) > 0 {
		return fmt.Errorf("failed to expand settings record: %v", errs)
	}
	return nil
}

func (s *Settings) init(ctx AppContext) error {
	record, err := ctx.App.FindFirstRecordByFilter(
		schema.CollectionSettings,
		"",
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if record == nil {
		ctx.App.Logger().Warn("Settings record not found, create new one")
		record = core.NewRecord(GameCollections.Get(schema.CollectionSettings))
		record.Set(schema.SettingsSchema.BlockAllActions, true)
	} else {
		err = expandSettingsWithSeason(ctx, record)
		if err != nil {
			return err
		}
	}

	s.SetProxyRecord(record)

	return nil
}

func (s *Settings) EventEnded() bool {
	return s.GetBool(schema.SettingsSchema.EventEnded)
}

func (s *Settings) CurrentSeason() string {
	return s.GetString(schema.SettingsSchema.CurrentSeason)
}

func (s *Settings) CurrentSeasonDateStart() types.DateTime {
	return s.ExpandedOne(schema.SettingsSchema.CurrentSeason).GetDateTime(schema.SeasonSchema.SeasonDateStart)
}

func (s *Settings) CurrentSeasonDateEnd() types.DateTime {
	return s.ExpandedOne(schema.SettingsSchema.CurrentSeason).GetDateTime(schema.SeasonSchema.SeasonDateEnd)
}

func (s *Settings) CurrentWeek() int {
	return s.GetInt(schema.SettingsSchema.CurrentWeek)
}

func (s *Settings) SetCurrentWeek(w int) {
	s.Set(schema.SettingsSchema.CurrentWeek, w)
}

func (s *Settings) BlockAllActions() bool {
	return s.GetBool(schema.SettingsSchema.BlockAllActions)
}

func (s *Settings) MaxInventorySlots() int {
	return s.GetInt(schema.SettingsSchema.MaxInventorySlots)
}

func (s *Settings) PointsForDrop() int {
	return s.GetInt(schema.SettingsSchema.PointsForDrop)
}

func (s *Settings) DropsToJail() int {
	return s.GetInt(schema.SettingsSchema.DropsToJail)
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

func (s *Settings) DisableRefreshHltbTime() bool {
	return s.GetBool(schema.SettingsSchema.DisableRefreshHltbTime)
}

func (s *Settings) OnKillParser() *event.Hook[*OnKillParserEvent] {
	return s.onKillParser
}

func (s *Settings) IgdbForceUpdateGames() bool {
	return s.GetBool(schema.SettingsSchema.IgdbForceUpdateGames)
}
