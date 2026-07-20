package adventuria

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/event_stats"
	"adventuria/internal/adventuria/inventories"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/players"
	"adventuria/internal/adventuria/scope"
	"adventuria/internal/adventuria/settings"
	"adventuria/internal/adventuria/worlds"
	"adventuria/pkg/event"
	"adventuria/pkg/locker"
	"adventuria/pkg/pbtransaction"
	"context"
	"errors"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

type Game struct {
	pb            *pocketbase.PocketBase
	settings      *settings.Settings
	players       *players.Players
	cells         *cells.Cells
	actionEvents  *action_events.ActionEvents
	actions       *actions.Actions
	inventories   *inventories.Inventories
	effects       *effects.Effects
	worlds        *worlds.Worlds
	eventStats    *event_stats.EventStats
	playersLocker *locker.Locker[string]

	onKillParserEvent *event.Hook[*onKillParserEvent]
}

func Start(fn func(game *Game, se *core.ServeEvent) error) (*Game, error) {
	g := &Game{
		pb:            pocketbase.New(),
		playersLocker: locker.New[string](),

		onKillParserEvent: &event.Hook[*onKillParserEvent]{},
	}

	migratecmd.MustRegister(g.pb, g.pb.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	g.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		if err := g.init(e.App); err != nil {
			return err
		}
		return e.Next()
	})

	g.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		return fn(g, e)
	})

	err := g.pb.Start()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) initScope(ctx context.Context, player *model.Player) (*scope.Scope, error) {
	s := scope.New(player)

	invs, err := g.inventories.GetAllByPlayerID(ctx, player.ID())
	if err != nil {
		return nil, err
	}

	err = g.effects.SubscribeActiveEffects(ctx, s.Events(), s.Player(), invs)
	if err != nil {
		return nil, err
	}

	err = g.effects.SubscribePersistentEffects(ctx, s.Events(), s.Player())
	if err != nil {
		return nil, err
	}

	err = g.actionEvents.ListenToActionEvents(s.Events(), s.Player())
	if err != nil {
		return nil, err
	}

	err = g.worlds.SubscribeEffects(ctx, s.Events(), s.Player(), s.Player().Progress().CurrentWorld())
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (g *Game) DoAction(
	ctx context.Context,
	pb core.App,
	playerId string,
	actionType model.ActionType,
	req model.ActionRequest,
) (any, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return nil, err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return nil, errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	s, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	if ok := g.actions.CanDo(ctx, s.Events(), s.Player(), actionType); !ok {
		return nil, errors.New("action is not available")
	}

	var res any
	err = pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		res, err = g.actions.Do(ctx, s.Events(), s.Player(), req, actionType)
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

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	s, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	canUse, err := g.inventories.CanUseItem(ctx, s.Events(), s.Player(), itemId)
	if err != nil {
		return err
	}
	if !canUse {
		return errors.New("can't use item")
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		err = g.inventories.UseItem(ctx, s.Events(), s.Player(), itemId)
		if err != nil {
			return err
		}

		err = s.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
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

	player, err := g.players.GetByID(ctx, playerId, currentSeason)
	if err != nil {
		return err
	}

	if ok := g.playersLocker.TryLock(playerId); !ok {
		return errs.ErrPlayerIsBusy
	}
	defer g.playersLocker.Unlock(playerId)

	s, err := g.initScope(ctx, player)
	if err != nil {
		return err
	}

	canDrop, err := g.inventories.CanDropItem(ctx, playerId, itemId)
	if err != nil {
		return err
	}
	if !canDrop {
		return errors.New("can't drop item")
	}

	item, err := g.inventories.GetPlayerInventoryItemByID(ctx, playerId, itemId)
	if err != nil {
		return err
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		err = g.inventories.DropItem(ctx, s.Events(), s.Player(), item)
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

	s, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	availableActions := g.actions.AvailableActions(ctx, s.Events(), s.Player())

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

	s, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	return g.effects.GetView(ctx, s.Events(), s.Player(), effectId)
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

	s, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	return g.actions.GetView(ctx, s.Events(), s.Player(), actionType)
}

func (g *Game) EventStats(ctx context.Context) (*event_stats.EventStatsData, error) {
	currentSeason, err := g.settings.CurrentSeason(ctx)
	if err != nil {
		return nil, err
	}

	return g.eventStats.ComputeStats(ctx, currentSeason)
}

func (g *Game) IsActionsBlocked(ctx context.Context) (bool, error) {
	return g.settings.IsActionsBlocked(ctx)
}

func (g *Game) CurrentSeason(ctx context.Context) (string, error) {
	return g.settings.CurrentSeason(ctx)
}

func (g *Game) IsEventEnded(ctx context.Context) (bool, error) {
	return g.settings.IsEventEnded(ctx)
}
