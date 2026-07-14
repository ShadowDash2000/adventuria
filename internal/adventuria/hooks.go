package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

func (g *Game) bindHooks(pb core.App) {
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
		ctx := context.Background()

		settings, err := g.settings.GetFirstOrDefault(ctx)
		if err != nil {
			return err
		}

		player, err := g.players.GetByID(ctx, e.Record.GetString(schema.InventorySchema.Player), settings.CurrentSeason())
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
}
