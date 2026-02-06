package adventuria

import (
	"adventuria/pkg/collections"
	"adventuria/pkg/event"
	"context"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/types"
)

var (
	PocketBase      core.App
	GameUsers       *Users
	GameCells       *Cells
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
	GameUsers = NewUsers(ctx)
	GameActions = NewActions()
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
	_ = NewTimers(ctx)

	return nil
}

func (g *Game) DoAction(app core.App, userId string, actionType ActionType, req ActionRequest) (*ActionResult, error) {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return &ActionResult{
			Success: false,
			Error:   "request error: user not found",
		}, nil
	}

	if user.isInAction() {
		return &ActionResult{
			Success: false,
			Error:   "request error: user is already in action",
		}, nil
	}

	user.setIsInAction(true)

	if ok := GameActions.CanDo(ctx, user, actionType); !ok {
		return &ActionResult{
			Success: false,
			Error:   "request error: cannot do action",
		}, nil
	}

	var (
		res    *ActionResult
		txUser User
	)
	err = ctx.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}

		txUser, err = NewUser(ctx, userId)
		if err != nil {
			return err
		}

		txUser.setIsInAction(true)

		res, err = GameActions.Do(ctx, txUser, req, actionType)
		if err != nil {
			app.Logger().Error(
				"Failed to complete user action",
				"error", err,
			)
			return err
		} else if res.Error != "" {
			return errors.New(res.Error)
		}

		return nil
	})
	if err != nil {
		user.setIsInAction(false)
		if res != nil {
			return res, err
		}
		return &ActionResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	user.Close(ctx)
	GameUsers.Update(txUser)
	defer txUser.setIsInAction(false)

	eventRes, err := txUser.OnAfterAction().Trigger(&OnAfterActionEvent{
		AppContext: ctx,
		ActionType: actionType,
	})
	if eventRes != nil && !eventRes.Success {
		return &ActionResult{
			Success: false,
			Error:   eventRes.Error,
		}, fmt.Errorf("doAction(): %w", err)
	}
	if err != nil {
		app.Logger().Error("Failed to trigger onAfterActionEvent", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	err = app.Save(txUser.LastAction().ProxyRecord())
	if err != nil {
		app.Logger().Error("Failed to save latest user action", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	err = app.Save(txUser.ProxyRecord())
	if err != nil {
		app.Logger().Error("Failed to save user", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	return res, nil
}

type UseItemRequest struct {
	InvItemId string         `json:"itemId"`
	Data      map[string]any `json:"data"`
}

func (g *Game) UseItem(app core.App, userId string, req UseItemRequest) error {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.setIsInAction(true)

	var (
		eventRes *event.Result
		txUser   User
	)
	err = ctx.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}

		txUser, err = NewUser(ctx, userId)
		if err != nil {
			return err
		}

		txUser.setIsInAction(true)

		onUseSuccess, onUseFail, err := txUser.Inventory().UseItem(ctx, req.InvItemId)
		if err != nil {
			return err
		}

		eventRes, err = txUser.OnAfterItemUse().Trigger(&OnAfterItemUseEvent{
			AppContext: ctx,
			InvItemId:  req.InvItemId,
			Data:       req.Data,
		})
		if eventRes != nil && !eventRes.Success {
			onUseFail()
			return errors.New(eventRes.Error)
		}
		if err != nil {
			onUseFail()
			return err
		}

		return onUseSuccess()
	})
	if err != nil {
		user.setIsInAction(false)
		return err
	}

	user.Close(ctx)
	GameUsers.Update(txUser)
	defer txUser.setIsInAction(false)

	eventRes, err = txUser.OnAfterAction().Trigger(&OnAfterActionEvent{
		AppContext: ctx,
		ActionType: "useItem",
	})
	if eventRes != nil && !eventRes.Success {
		return errors.New(eventRes.Error)
	}
	if err != nil {
		return err
	}

	err = app.Save(txUser.LastAction().ProxyRecord())
	if err != nil {
		return err
	}

	err = app.Save(txUser.ProxyRecord())
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) DropItem(app core.App, userId, itemId string) error {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	item, ok := user.Inventory().GetItemById(itemId)
	if !ok {
		return errors.New("item not found in inventory")
	}

	err = user.Inventory().DropItem(ctx, itemId)
	if err != nil {
		return err
	}

	if itemPrice := item.Price(); itemPrice > 0 {
		user.SetBalance(user.Balance() + itemPrice/2)

		err = app.Save(user.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) StartTimer(app core.App, userId string) error {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	return user.Timer().Start(ctx)
}

func (g *Game) StopTimer(app core.App, userId string) error {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	return user.Timer().Stop(ctx)
}

func (g *Game) GetTimeLeft(app core.App, userId string) (int64, bool, types.DateTime, error) {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer().GetTimeLeft(), user.Timer().IsActive(), GameSettings.NextTimerResetDate(), nil
}

func (g *Game) GetAvailableActions(app core.App, userId string) ([]ActionType, error) {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	var actions []ActionType
	for t := range GameActions.AvailableActions(ctx, user) {
		actions = append(actions, t)
	}

	return actions, nil
}

func (g *Game) Context() context.Context {
	return g.ctx
}

func (g *Game) GetItemEffectVariants(app core.App, userId, invItemId, effectId string) (any, error) {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	invItem, ok := user.Inventory().GetItemById(invItemId)
	if !ok {
		return nil, errors.New("inventory item not found")
	}

	return invItem.GetEffectVariants(ctx, effectId)
}

func (g *Game) GetActionVariants(app core.App, userId, actionType string) (any, error) {
	ctx := AppContext{App: app}
	user, err := GameUsers.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if ok := GameActions.CanDo(ctx, user, ActionType(actionType)); !ok {
		return &ActionResult{
			Success: false,
			Error:   "request error: cannot do action",
		}, nil
	}

	variants := GameActions.GetVariants(ctx, user, ActionType(actionType))
	if variants == nil {
		return nil, errors.New("action variants not found")
	}

	return variants, nil
}
