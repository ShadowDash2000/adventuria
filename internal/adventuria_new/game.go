package adventuria_new

import (
	"adventuria/internal/adventuria_new/actions"
	customActions "adventuria/internal/adventuria_new/actions/custom"
	actionsRepo "adventuria/internal/adventuria_new/actions/repository"
	"adventuria/internal/adventuria_new/activities"
	activitiesRepo "adventuria/internal/adventuria_new/activities/repository"
	"adventuria/internal/adventuria_new/activity_filters"
	activityFiltersRepo "adventuria/internal/adventuria_new/activity_filters/repository"
	"adventuria/internal/adventuria_new/board"
	"adventuria/internal/adventuria_new/cells"
	customCells "adventuria/internal/adventuria_new/cells/custom"
	cellsRepo "adventuria/internal/adventuria_new/cells/repository"
	"adventuria/internal/adventuria_new/effects"
	customEffects "adventuria/internal/adventuria_new/effects/custom"
	effectsRepo "adventuria/internal/adventuria_new/effects/repository"
	"adventuria/internal/adventuria_new/genres"
	genresRepo "adventuria/internal/adventuria_new/genres/repository"
	"adventuria/internal/adventuria_new/inventories"
	inventoriesRepo "adventuria/internal/adventuria_new/inventories/repository"
	"adventuria/internal/adventuria_new/items"
	itemsRepo "adventuria/internal/adventuria_new/items/repository"
	"adventuria/internal/adventuria_new/model"
	"adventuria/internal/adventuria_new/player_progress"
	progressRepo "adventuria/internal/adventuria_new/player_progress/repository"
	"adventuria/internal/adventuria_new/players"
	playersRepo "adventuria/internal/adventuria_new/players/repository"
	"adventuria/internal/adventuria_new/reviews"
	reviewsRepo "adventuria/internal/adventuria_new/reviews/repository"
	"adventuria/internal/adventuria_new/scope"
	"adventuria/internal/adventuria_new/seasons"
	seasonsRepo "adventuria/internal/adventuria_new/seasons/repository"
	"adventuria/internal/adventuria_new/settings"
	settingsRepo "adventuria/internal/adventuria_new/settings/repository"
	"adventuria/internal/adventuria_new/worlds"
	worldsRepo "adventuria/internal/adventuria_new/worlds/repository"
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
	seasonsRepository := seasonsRepo.NewRepository(pb)
	settingsRepository := settingsRepo.NewRepository(pb)
	cellsRepository := cellsRepo.NewRepository(pb)
	worldsRepository := worldsRepo.NewRepository(pb)
	actionsRepository := actionsRepo.NewRepository(pb)
	progressRepository := progressRepo.NewRepository(pb)
	playersRepository := playersRepo.NewRepository(pb)
	inventoriesRepository := inventoriesRepo.NewRepository(pb)
	effectsRepository := effectsRepo.NewRepository(pb)
	activitiesRepository := activitiesRepo.NewRepository(pb)
	activityFiltersRepository := activityFiltersRepo.NewRepository(pb)
	itemsRepository := itemsRepo.NewRepository(pb)
	genresRepository := genresRepo.NewRepository(pb)
	reviewsRepository := reviewsRepo.NewRepository(pb)

	seasonsService := seasons.NewSeasons(seasonsRepository)
	settingsService := settings.NewSettings(settingsRepository, seasonsService)
	worldsService := worlds.NewWorlds(worldsRepository)
	cellsService := cells.NewCells(cellsRepository)
	actionsService := actions.NewActions(actionsRepository, worldsService, cellsService)
	progressService := player_progress.NewPlayerProgress(progressRepository, worldsService)
	playersService := players.NewPlayers(
		playersRepository,
		actionsService,
		progressService,
		seasonsService,
	)
	effectsService := effects.NewEffects(effectsRepository, inventoriesRepository)
	inventoriesService := inventories.NewInventories(inventoriesRepository, effectsService, itemsRepository)
	activitiesService := activities.NewActivities(activitiesRepository)
	activityFiltersService := activity_filters.NewActivityFilters(activityFiltersRepository)
	itemsService := items.NewItems(itemsRepository)
	boardService := board.NewBoard(actionsService, progressService, cellsService, worldsService)
	genresService := genres.NewGenres(genresRepository)
	reviewsService := reviews.NewReviews(reviewsRepository)

	g.settings = settingsService
	g.players = playersService
	g.cells = cellsService
	g.actions = actionsService
	g.inventories = inventoriesService
	g.effects = effectsService

	customCells.RegisterCells(
		activitiesService,
		activityFiltersService,
		itemsService,
		cellsService,
		actionsService,
		boardService,
	)

	customEffects.RegisterEffects(
		actionsService,
		cellsService,
		genresService,
		activityFiltersService,
	)

	customActions.RegisterActions(
		cellsService,
		reviewsService,
		playersService,
		settingsService,
		boardService,
		actionsService,
		itemsService,
		inventoriesService,
	)

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
		return nil, errors.New("player is busy")
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
		return errors.New("player is busy")
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

		err = s.Events().OnAfterItemUse().Trigger(&model.OnAfterItemUseEvent{
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
		return errors.New("player is busy")
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

func (g *Game) GetActionView(
	ctx context.Context,
	playerId string,
	actionType model.ActionType,
) (any, error) {
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
