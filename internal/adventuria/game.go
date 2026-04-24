package adventuria

import (
	"adventuria/pkg/collections"
	"adventuria/pkg/result"
	"context"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

var (
	PocketBase      core.App
	GamePlayers     *Players
	GameCells       *Cells
	GameWorlds      *Worlds
	GameItems       *Items
	GameCollections *collections.Collections
	GameSettings    *Settings
	GameActions     *Actions
)

type AppContext struct {
	App core.App
}

type Game struct {
	pb     *pocketbase.PocketBase
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Game {
	return &Game{}
}

func (g *Game) Start(fn func(se *core.ServeEvent) error) error {
	g.pb = pocketbase.New()
	g.ctx, g.cancel = context.WithCancel(context.Background())
	PocketBase = g.pb

	migratecmd.MustRegister(g.pb, g.pb.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	g.pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		if err := g.init(AppContext{App: e.App}); err != nil {
			return err
		}
		return e.Next()
	})

	g.pb.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		g.cancel()
		return e.Next()
	})

	g.pb.OnServe().BindFunc(fn)

	return g.pb.Start()
}

func (g *Game) init(ctx AppContext) error {
	var err error

	GameCollections = collections.NewCollections(PocketBase)
	GamePlayers = NewPlayers(ctx)
	GameActions = NewActions(ctx)
	GameWorlds, err = NewWorlds(ctx)
	if err != nil {
		return err
	}
	GameCells, err = NewCells(ctx)
	if err != nil {
		return err
	}
	GameItems, err = NewItems(ctx)
	if err != nil {
		return err
	}
	GameSettings, err = NewSettings(ctx)
	if err != nil {
		return err
	}

	_ = NewInventories(ctx)
	_ = NewEffectVerifier(ctx)
	_ = NewCellVerifier(ctx)

	BindActivitiesHooks(ctx)

	return nil
}

func (g *Game) DoAction(app core.App, playerId string, actionType ActionType, req ActionRequest) (*result.Result, error) {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return result.Err("player not found"), nil
	}

	if player.Locked() {
		return result.Err("player is already in action"), nil
	}

	player.Lock()
	defer player.Unlock()

	if ok := GameActions.CanDo(ctx, player, actionType); !ok {
		return result.Err("action is not available"), nil
	}

	var res *result.Result
	err = ctx.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}

		res, err = GameActions.Do(ctx, player, req, actionType)
		if err != nil {
			txApp.Logger().Error(
				"Failed to complete player action",
				"error", err,
			)
			return err
		} else if res.Error != "" {
			return errors.New(res.Error)
		}

		_, err = player.OnAfterAction().Trigger(&OnAfterActionEvent{
			AppContext: ctx,
			ActionType: actionType,
		})
		if err != nil {
			return err
		}

		err = txApp.Save(player.LastAction().ProxyRecord())
		if err != nil {
			txApp.Logger().Error("Failed to save latest player action", "error", err)
			res = result.Err(err.Error())
			return err
		}

		err = txApp.Save(player.Progress().ProxyRecord())
		if err != nil {
			txApp.Logger().Error("Failed to save player", "error", err)
			return err
		}

		return nil
	})
	if err != nil {
		_ = player.Refetch(ctx)
	}

	return res, err
}

type UseItemRequest struct {
	InvItemId string         `json:"itemId"`
	Data      map[string]any `json:"data"`
}

func (g *Game) UseItem(app core.App, playerId string, req UseItemRequest) (*result.Result, error) {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return result.Err("player not found"), err
	}

	player.Lock()
	defer player.Unlock()

	if ok := player.Inventory().CanUseItem(ctx, req.InvItemId); !ok {
		return result.Err(fmt.Sprintf("item %s cannot be used", req.InvItemId)), nil
	}

	var res *result.Result
	err = ctx.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}

		onUseSuccess, onUseFail, err := player.Inventory().UseItem(ctx, req.InvItemId)
		if err != nil {
			return err
		}

		res, err = player.OnAfterItemUse().Trigger(&OnAfterItemUseEvent{
			AppContext: ctx,
			InvItemId:  req.InvItemId,
			Data:       req.Data,
		})
		if err != nil {
			onUseFail()
			return err
		}
		if res.Failed() {
			onUseFail()
			return errors.New(res.Error)
		}

		err = onUseSuccess()
		if err != nil {
			return err
		}

		res, err = player.OnAfterAction().Trigger(&OnAfterActionEvent{
			AppContext: ctx,
			ActionType: "useItem",
		})
		if err != nil {
			return err
		}
		if res.Failed() {
			return errors.New(res.Error)
		}

		err = txApp.Save(player.LastAction().ProxyRecord())
		if err != nil {
			return err
		}

		err = txApp.Save(player.Progress().ProxyRecord())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		_ = player.Refetch(ctx)
		return res, err
	}

	return res, nil
}

func (g *Game) DropItem(app core.App, playerId, itemId string) error {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return err
	}

	player.Lock()
	defer player.Unlock()

	item, ok := player.Inventory().GetItemById(itemId)
	if !ok {
		return errors.New("item not found in inventory")
	}

	err = player.Inventory().DropItem(ctx, itemId)
	if err != nil {
		return err
	}

	if itemPrice := item.Price(); itemPrice > 0 {
		player.Progress().AddBalance(itemPrice / 2)
		err = app.Save(player.Progress().ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) GetAvailableActions(app core.App, playerId string) ([]ActionType, error) {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return nil, err
	}

	var actions []ActionType
	for t := range GameActions.AvailableActions(ctx, player) {
		actions = append(actions, t)
	}

	return actions, nil
}

func (g *Game) Context() context.Context {
	return g.ctx
}

func (g *Game) GetItemEffectVariants(app core.App, playerId, invItemId, effectId string) (any, error) {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return nil, err
	}

	invItem, ok := player.Inventory().GetItemById(invItemId)
	if !ok {
		return nil, errors.New("inventory item not found")
	}

	return invItem.GetEffectVariants(ctx, effectId)
}

func (g *Game) GetActionVariants(app core.App, playerId, actionType string) (*result.Result, error) {
	ctx := AppContext{App: app}
	player, err := GamePlayers.GetByID(ctx, playerId)
	if err != nil {
		return result.Err("player not found"), err
	}

	if ok := GameActions.CanDo(ctx, player, ActionType(actionType)); !ok {
		return result.Err("action is not available"), nil
	}

	variants := GameActions.GetVariants(ctx, player, ActionType(actionType))
	if variants == nil {
		return result.Err("action variants not found"), nil
	}

	return result.Ok().WithData(variants), nil
}
