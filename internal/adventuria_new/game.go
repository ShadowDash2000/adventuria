package adventuria_new

import (
	"adventuria/internal/adventuria_new/actions"
	customActions "adventuria/internal/adventuria_new/actions/custom"
	"adventuria/internal/adventuria_new/activities"
	"adventuria/internal/adventuria_new/cells"
	customCells "adventuria/internal/adventuria_new/cells/custom"
	"adventuria/internal/adventuria_new/effects"
	customEffects "adventuria/internal/adventuria_new/effects/custom"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/inventories"
	"adventuria/internal/adventuria_new/model"
	customOutboxes "adventuria/internal/adventuria_new/outboxes/custom"
	"adventuria/internal/adventuria_new/players"
	"adventuria/internal/adventuria_new/scope"
	"adventuria/internal/adventuria_new/settings"
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
	actions       *actions.Actions
	inventories   *inventories.Inventories
	effects       *effects.Effects
	playersLocker *locker.Locker[string]
}

func Start(fn func(se *core.ServeEvent) error) (*Game, error) {
	g := &Game{
		pb:            pocketbase.New(),
		playersLocker: locker.New[string](),
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

	g.pb.OnServe().BindFunc(fn)

	err := g.pb.Start()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) init(pb core.App) error {
	registry := NewRegistry(pb)

	g.settings = registry.Settings()
	g.players = registry.Players()
	g.cells = registry.Cells()
	g.actions = registry.Actions()
	g.inventories = registry.Inventories()
	g.effects = registry.Effects()

	customCells.RegisterCells(
		registry.Activities(),
		registry.ActivityFilters(),
		registry.Items(),
		registry.Cells(),
		registry.Actions(),
		registry.Board(),
	)

	customEffects.RegisterEffects(
		registry.Actions(),
		registry.Cells(),
		registry.Genres(),
		registry.ActivityFilters(),
		registry.Inventories(),
		registry.Items(),
		registry.Activities(),
		registry.Players(),
		registry.Outboxes(),
		registry.Board(),
	)

	customEffects.RegisterPersistentEffects()

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
	registry.Outboxes().Start(context.Background())

	// hooks
	cells.BindHooks(pb)
	effects.BindHooks(pb)
	activities.BindHooks(pb, registry.RelationRepo())

	return nil
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

	return s, nil
}

func (g *Game) DoAction(
	ctx context.Context,
	pb core.App,
	playerId string,
	actionType model.ActionType,
	req model.ActionRequest,
) (any, error) {
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
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
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
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

	canUse, err := g.inventories.CanUseItem(ctx, s, itemId)
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
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
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

	item, err := g.inventories.GetPlayerInventoryItemByID(ctx, playerId, itemId)
	if err != nil {
		return err
	}

	return pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		return g.inventories.DropItem(ctx, s.Events(), s.Player(), item)
	})
}

func (g *Game) GetAvailableActions(ctx context.Context, playerId string) ([]model.ActionType, error) {
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
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
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
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
	settings, err := g.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	player, err := g.players.GetByID(ctx, playerId, settings.CurrentSeason())
	if err != nil {
		return nil, err
	}

	s, err := g.initScope(ctx, player)
	if err != nil {
		return nil, err
	}

	return g.actions.GetView(ctx, s.Events(), s.Player(), actionType)
}
