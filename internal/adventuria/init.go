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
		registry.Activities(),
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
	g.registerCrons(ctx, pb, registry)

	return nil
}
