package adventuria

import (
	"adventuria/internal/adventuria/action_events"
	customActionEvents "adventuria/internal/adventuria/action_events/custom"
	customActions "adventuria/internal/adventuria/actions/custom"
	"adventuria/internal/adventuria/activities"
	"adventuria/internal/adventuria/cells"
	customCells "adventuria/internal/adventuria/cells/custom"
	"adventuria/internal/adventuria/effects"
	customEffects "adventuria/internal/adventuria/effects/custom"
	customOutboxes "adventuria/internal/adventuria/outboxes/custom"
	"adventuria/internal/adventuria/stream_tracker"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

func (g *Game) init(ctx context.Context, pb core.App) error {
	registry := NewRegistry(pb, pb.Logger())

	g.settings = registry.Settings()
	g.players = registry.Players()
	g.cells = registry.Cells()
	g.cellEvents = registry.CellEventsSchedules()
	g.actions = registry.Actions()
	g.inventories = registry.Inventories()
	g.effects = registry.Effects()
	g.worlds = registry.Worlds()
	g.eventStats = registry.EventStats()

	customCells.RegisterCells(
		registry.Activities(),
		registry.ActivityFilters(),
		registry.Items(),
		registry.Cells(),
		registry.Actions(),
		registry.Board(),
	)

	customActionEvents.RegisterActionEvents(
		registry.Items(),
	)

	customEffects.RegisterEffects(
		registry.Actions(),
		registry.Cells(),
		registry.Genres(),
		registry.ActivityFilters(),
		registry.Inventories(),
		registry.Items(),
		registry.Activities(),
		registry.PlayerProgress(),
		registry.Outboxes(),
		registry.Board(),
	)

	customEffects.RegisterPersistentEffects(
		registry.ActivityFilters(),
	)

	customActions.RegisterActions(
		registry.Cells(),
		registry.Reviews(),
		registry.Players(),
		registry.Settings(),
		registry.Board(),
		registry.Actions(),
		registry.Items(),
		registry.Inventories(),
		registry.RollWheelRepo(),
	)

	customOutboxes.RegisterOutboxes(
		registry.PlayerProgress(),
	)

	// background tasks
	registry.Outboxes().Start(ctx)
	err := registry.StreamTracker().Start(ctx)
	if err != nil {
		return err
	}

	// hooks
	g.bindHooks(ctx, pb)
	cells.BindHooks(pb)
	effects.BindHooks(pb)
	action_events.BindHooks(pb)
	activities.BindHooks(pb, registry.RelationRepo())
	stream_tracker.BindHooks(pb, registry.StreamTracker())

	// crons
	pb.Cron().MustAdd("games_parser", "0 0 1 * *", func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		settings, err := g.settings.GetFirstOrDefault(ctx)
		if err != nil {
			return
		}

		unsub := g.onKillParser().BindFunc(func(ctx context.Context, e *onKillParserEvent) error {
			cancel()
			return e.Next()
		})
		defer unsub()

		if !settings.DisableSteamParser() {
			pb.Logger().Info("Started SteamSpy parser")

			err = registry.SteamSpy().Parse(ctx)
			if err != nil {
				pb.Logger().Error("SteamSpy parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished SteamSpy parser")
		}

		if !settings.DisableCheapsharkParser() {
			pb.Logger().Info("Started CheapShark parser")

			err = registry.CheapShark().Parse(ctx)
			if err != nil {
				pb.Logger().Error("CheapShark parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished CheapShark parser")
		}

		if !settings.DisableHltbParser() {
			pb.Logger().Info("Started HLTB parser")

			err = registry.HLTB().Parse(ctx)
			if err != nil {
				pb.Logger().Error("HLTB parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished HLTB parser")
		}

		if !settings.DisableIgdbParser() {
			err = registry.IGDB().ParsePlatforms(ctx, 500)
			if err != nil {
				return
			}
			err = registry.IGDB().ParseGenres(ctx, 500)
			if err != nil {
				return
			}
			err = registry.IGDB().ParseGameTypes(ctx, 500)
			if err != nil {
				return
			}
			if !settings.DisableIgdbGamesParser() {
				err = registry.IGDB().ParseGames(ctx, settings.IgdbFilter().Build(), 500)
				if err != nil {
					return
				}
			}
		}
	})
	pb.Cron().MustAdd("cell_events_scheduler", "*/1 * * * *", func() {
		err := g.cellEvents.CheckEventsSchedules(ctx)
		if err != nil {
			return
		}
	})

	return nil
}
