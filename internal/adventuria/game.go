package adventuria

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cell_events_schedules"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/inventories"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"adventuria/internal/adventuria/player_progress"
	"adventuria/internal/adventuria/players"
	"adventuria/internal/adventuria/scope"
	"adventuria/internal/adventuria/settings"
	"adventuria/internal/adventuria/worlds"
	"adventuria/pkg/event"
	"adventuria/pkg/locker"
	"adventuria/pkg/pbtransaction"
	"context"
	"errors"

	"ganigeorgiev.com/pocketbase/pocketbase"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/plugins/migratecmd"
)

type Game struct {
	pb               *pocketbase.PocketBase
	eventsSubscriber *scope.EventsSubscriber
	settings         *settings.Settings
	players          *players.Players
	playerProgress   *player_progress.PlayerProgress
	cells            *cells.Cells
	cellEvents       *cell_events_schedules.CellEventsSchedules
	actions          *actions.Actions
	inventories      *inventories.Inventories
	effects          *effects.Effects
	worlds           *worlds.Worlds
	board            *board.Board
	playersLocker    *locker.Locker[string]
	playerFactory    *player.Factory

	onKillParserEvent *event.Hook[*onKillParserEvent]
}

func Start(fn func(game *Game, registry *Registry, se *core.ServeEvent) error) (*Game, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := &Game{
		pb: pocketbase.New(),

		onKillParserEvent: &event.Hook[*onKillParserEvent]{},
	}

	migratecmd.MustRegister(g.pb, g.pb.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	g.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		registry, err := g.init(ctx, e.App)
		if err != nil {
			return err
		}

		err = fn(g, registry, e)
		if err != nil {
			return err
		}

		return e.Next()
	})

	err := g.pb.Start()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) initScope(ctx context.Context, player *model.Player) (*player.Actor, error) {
	actor := g.playerFactory.New(model.NewEvents(), player)

	inventory, err := g.inventories.GetAllWithEffectsByPlayerID(ctx, player.ID())
	if err != nil {
		return nil, err
	}

	err = g.eventsSubscriber.Subscribe(ctx, actor, inventory)
	if err != nil {
		return nil, err
	}

	return actor, nil
}

func (g *Game) DoAction(
	ctx context.Context,
	pb core.App,
	playerId string,
	actionType model.ActionType,
	req actions.Request,
) (any, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return nil, errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return nil, err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	if ok := g.actions.CanDo(ctx, actor, actionType); !ok {
		return nil, errors.New("action is not available")
	}

	var res any
	err = pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		res, err = g.actions.Do(ctx, actor, req, actionType)
		if err != nil {
			return err
		}

		err = g.players.Save(ctx, player)
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}

func (g *Game) AddItemByID(ctx context.Context, pb core.App, playerId, itemId string, source item_received.Source) error {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		_, err = g.inventories.AddItemByID(ctx, actor, itemId, source)
		if err != nil {
			return err
		}

		err = g.players.Save(ctx, player)
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *Game) UseItem(
	ctx context.Context,
	pb core.App,
	playerId string,
	itemId string,
	data map[string]any,
) error {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	canUse, err := g.inventories.CanUseItem(ctx, actor, itemId)
	if err != nil {
		return err
	}
	if !canUse {
		return errors.New("can't use item")
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		err = g.inventories.UseItem(ctx, actor, itemId)
		if err != nil {
			return err
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: itemId,
			Data:      data,
		})
		if err != nil {
			return err
		}

		err = g.players.Save(ctx, player)
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *Game) DropItem(ctx context.Context, pb core.App, playerId, itemId string) error {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	item, err := g.inventories.GetPlayerInventoryItemByID(ctx, playerId, itemId)
	if err != nil {
		return err
	}

	if !item.CanDrop() {
		return errors.New("can't drop item")
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		err = g.inventories.DropItem(ctx, actor.Events(), actor.Player(), item)
		if err != nil {
			return err
		}

		err = g.players.Save(ctx, player)
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *Game) MoveToCellID(ctx context.Context, pb core.App, playerId, cellId string) error {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		_, err = g.board.MoveToCellId(ctx, actor, cellId)
		if err != nil {
			return err
		}

		err = g.players.Save(ctx, player)
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *Game) GetAvailableActions(ctx context.Context, playerId string) ([]model.ActionType, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return nil, err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	availableActions := g.actions.AvailableActions(ctx, actor)

	return availableActions, nil
}

func (g *Game) GetEffectView(ctx context.Context, playerId, effectId string) (any, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return nil, err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	return g.effects.GetView(ctx, actor.Events(), actor.Player(), effectId)
}

func (g *Game) GetActionView(ctx context.Context, playerId string, actionType model.ActionType) (any, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return nil, err
	}

	actor, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	return g.actions.GetView(ctx, actor.Events(), actor.Player(), actionType)
}

func (g *Game) IsActionsBlocked(ctx context.Context) error {
	if g.cellEvents.IsRunning() {
		return errs.ErrCellEventSchedulerRunning
	}

	isBlocked, err := g.settings.IsActionsBlocked(ctx)
	if err != nil {
		return err
	}

	if isBlocked {
		return errs.ErrAllActionsBlocked
	}

	return nil
}
