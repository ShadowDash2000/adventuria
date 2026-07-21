package adventuria

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/schema"
	"context"
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
}
