package adventuria

import (
	"adventuria/internal/adventuria/cell_events_schedules"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/player_events/custom/energy_changed"
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria/settings"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbrealtime"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func (g *Game) bindHooks(ctx context.Context, pb core.App) {
	pb.OnRecordUpdate(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		if ok := e.Record.GetBool(schema.SettingsSchema.KillParser); ok {
			err := g.onKillParser().Trigger(e.Context, &onKillParserEvent{})
			if err != nil {
				e.App.Logger().Error("Failed to trigger kill parser event", "err", err)
			}
			e.Record.Set(schema.SettingsSchema.KillParser, false)
		}
		return e.Next()
	})

	pb.OnRecordEnrich(schema.CollectionInventory).BindFunc(func(e *core.RecordEnrichEvent) error {
		currentSeason, err := g.settings.CurrentSeason(ctx)
		if err != nil {
			return err
		}

		player, err := g.players.GetByID(ctx, e.Record.GetString(schema.InventorySchema.Player), currentSeason)
		if err != nil {
			return err
		}

		s, err := g.initScope(ctx, player)
		if err != nil {
			return err
		}

		canUse, err := g.inventories.CanUseItem(ctx, s.Events(), s.Player(), e.Record.Id)
		if err != nil {
			return err
		}
		canDrop, err := g.inventories.CanDropItem(ctx, player.ID(), e.Record.Id)
		if err != nil {
			return err
		}

		e.Record.WithCustomData(true)
		e.Record.Set("can_use", canUse)
		e.Record.Set("can_drop", canDrop)

		return e.Next()
	})

	pb.OnRecordEnrich(schema.CollectionCells).BindFunc(func(e *core.RecordEnrichEvent) error {
		cellEventId, err := g.cellEvents.GetIDByActiveCellID(ctx, e.Record.Id)
		if err != nil {
			if errors.Is(err, errs.ErrCellEventScheduleNotFound) {
				return e.Next()
			}
			return err
		}

		var record core.Record
		err = e.App.RecordQuery(schema.CollectionCellEventsSchedule).
			WithContext(ctx).
			Where(dbx.HashExp{
				schema.CellEventsScheduleSchema.Id: cellEventId,
			}).
			Limit(1).
			One(&record)
		if err != nil {
			return err
		}

		e.Record.WithCustomData(true)
		e.Record.Set("cell_event", record)

		return e.Next()
	})

	pb.OnRecordEnrich(schema.CollectionActions).BindFunc(func(e *core.RecordEnrichEvent) error {
		var records []*core.Record
		err := e.App.RecordQuery(schema.CollectionPlayerEvents).
			WithContext(ctx).
			Where(dbx.HashExp{
				schema.PlayerEventsSchema.Action: e.Record.Id,
			}).
			OrderBy("created DESC").
			Limit(20).
			All(&records)
		if err != nil {
			return err
		}

		e.Record.WithCustomData(true)
		e.Record.Set("player_events", records)

		return e.Next()
	})

	pb.OnRecordAfterDeleteSuccess(schema.CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		reviewId := e.Record.GetString(schema.ActionSchema.Review)
		if reviewId == "" {
			return e.Next()
		}

		var record core.Record
		err := e.App.RecordQuery(schema.CollectionReviews).
			WithContext(ctx).
			Where(dbx.HashExp{
				schema.ReviewSchema.Id: reviewId,
			}).
			Limit(1).
			One(&record)
		if err != nil {
			return err
		}

		err = e.App.DeleteWithContext(ctx, &record)
		if err != nil {
			return err
		}

		return e.Next()
	})

	pb.OnRecordDeleteExecute(schema.CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		return e.App.RunInTransaction(func(txApp core.App) error {
			var records []*core.Record
			err := txApp.RecordQuery(schema.CollectionPlayerEvents).
				WithContext(ctx).
				Where(dbx.HashExp{
					schema.PlayerEventsSchema.Action: e.Record.Id,
				}).
				All(&records)
			if err != nil {
				return err
			}

			for _, record := range records {
				err = txApp.DeleteWithContext(ctx, record)
				if err != nil {
					return err
				}
			}

			e.App = txApp
			return e.Next()
		})
	})

	pb.OnRecordAfterDeleteSuccess(schema.CollectionReviews).BindFunc(func(e *core.RecordEvent) error {
		filesIds := e.Record.GetStringSlice(schema.ReviewSchema.Files)
		if len(filesIds) == 0 {
			return e.Next()
		}

		var records []*core.Record
		err := e.App.RecordQuery(schema.CollectionReviewsFiles).
			WithContext(ctx).
			Where(dbx.In(
				schema.ReviewsFilesSchema.Id,
				pbhelper.SliceToAny(filesIds)...,
			)).
			All(&records)
		if err != nil {
			return err
		}

		for _, record := range records {
			err = e.App.DeleteWithContext(ctx, record)
			if err != nil {
				return err
			}
		}

		return e.Next()
	})

	pb.OnRecordAfterDeleteSuccess(schema.CollectionPlayersProgress).BindFunc(func(e *core.RecordEvent) error {
		playerId := e.Record.GetString(schema.PlayerProgressSchema.Player)
		seasonId := e.Record.GetString(schema.PlayerProgressSchema.Season)

		var record core.Record
		err := e.App.RecordQuery(schema.CollectionPlayerStats).
			WithContext(ctx).
			Where(dbx.HashExp{
				schema.PlayerStatsSchema.Player: playerId,
				schema.PlayerStatsSchema.Season: seasonId,
			}).
			Limit(1).
			One(&record)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.Next()
			}
			return err
		}

		err = e.App.DeleteWithContext(ctx, &record)
		if err != nil {
			return err
		}

		return e.Next()
	})

	g.cellEvents.OnUpdateStart().BindFunc(func(ctx context.Context, e *cell_events_schedules.UpdateStartEvent) error {
		err := pbrealtime.Notify(ctx, pb, "cell_events_scheduler_start", nil)
		if err != nil {
			pb.Logger().Error("Failed to notify cell events scheduler start", "error", err)
		}

		return e.Next()
	})

	g.cellEvents.OnUpdateEnd().BindFunc(func(ctx context.Context, e *cell_events_schedules.UpdateEndEvent) error {
		err := pbrealtime.Notify(ctx, pb, "cell_events_scheduler_end", nil)
		if err != nil {
			pb.Logger().Error("Failed to notify cell events scheduler end", "error", err)
		}

		return e.Next()
	})

	g.settings.OnNewDay().BindFunc(func(ctx context.Context, e *settings.OnNewDayEvent) error {
		progresses, err := g.playerProgress.GetAllBySeasonID(ctx, e.Season.ID())
		if err != nil {
			return err
		}

		energyAmount := e.Settings.EnergyPerDay() * e.DaysPassed
		for _, p := range progresses {
			_, err = g.playerProgress.ChangeEnergyOutbox(ctx, p.ID(), energyAmount, energy_changed.SourceDailyReward)
			if err != nil {
				return err
			}
		}

		return e.Next()
	})

	pb.OnRecordAfterCreateSuccess(schema.CollectionPlayersProgress).BindFunc(func(e *core.RecordEvent) error {
		currentSeason, err := g.settings.CurrentSeason(ctx)
		if err != nil {
			return err
		}

		if e.Record.GetString(schema.PlayerProgressSchema.Season) == currentSeason {
			err = pbrealtime.Notify(ctx, e.App, "new_player", nil)
			if err != nil {
				e.App.Logger().Error("Failed to notify new player", "error", err)
			}
		}

		return e.Next()
	})
}
