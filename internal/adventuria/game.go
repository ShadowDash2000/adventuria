package adventuria

import (
	"adventuria/pkg/collections"
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

type Game struct {
	pb     *pocketbase.PocketBase
	ctx    context.Context
	cancel context.CancelFunc
	ef     *EffectVerifier
	cf     *CellVerifier
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
		if err := g.init(); err != nil {
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

func (g *Game) init() error {
	var err error

	GameCollections = collections.NewCollections(PocketBase)
	GameUsers = NewUsers()
	GameActions = NewActions()
	GameCells, err = NewCells()
	if err != nil {
		return err
	}
	GameItems, err = NewItems()
	if err != nil {
		return err
	}
	GameSettings, err = NewSettings()
	if err != nil {
		return err
	}

	g.ef = NewEffectVerifier()
	g.cf = NewCellVerifier()
	return nil
}

func (g *Game) DoAction(actionType ActionType, userId string, req ActionRequest) (*ActionResult, error) {
	user, err := GameUsers.GetByID(userId)
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
	defer user.setIsInAction(false)

	if ok := GameActions.CanDo(user, actionType); !ok {
		return &ActionResult{
			Success: false,
			Error:   "request error: cannot do action",
		}, nil
	}

	res, err := GameActions.Do(user, req, actionType)
	if err != nil {
		PocketBase.Logger().Error("Failed to complete user action", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	} else if res.Error != "" {
		return res, nil
	}

	eventRes, err := user.OnAfterAction().Trigger(&OnAfterActionEvent{
		ActionType: actionType,
	})
	if eventRes != nil && !eventRes.Success {
		return &ActionResult{
			Success: false,
			Error:   eventRes.Error,
		}, fmt.Errorf("doAction(): %w", err)
	}
	if err != nil {
		PocketBase.Logger().Error("Failed to trigger onAfterActionEvent", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	err = PocketBase.Save(user.LastAction().ProxyRecord())
	if err != nil {
		PocketBase.Logger().Error("Failed to save latest user action", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	err = PocketBase.Save(user.ProxyRecord())
	if err != nil {
		PocketBase.Logger().Error("Failed to save user", "error", err)
		return &ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("doAction(): %w", err)
	}

	return res, nil
}

type UseItemRequest map[string]any

func (g *Game) UseItem(userId, itemId string, req UseItemRequest) error {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return err
	}

	onUseSuccess, onUseFail, err := user.Inventory().UseItem(itemId)
	if err != nil {
		return err
	}

	eventRes, err := user.OnAfterItemUse().Trigger(&OnAfterItemUseEvent{
		InvItemId: itemId,
		Request:   req,
	})
	if eventRes != nil && !eventRes.Success {
		onUseFail()
		return errors.New(eventRes.Error)
	}
	if err != nil {
		onUseFail()
		return err
	}

	if err = onUseSuccess(); err != nil {
		return err
	}

	eventRes, err = user.OnAfterAction().Trigger(&OnAfterActionEvent{
		ActionType: "useItem",
	})
	if eventRes != nil && !eventRes.Success {
		return errors.New(eventRes.Error)
	}
	if err != nil {
		return err
	}

	err = PocketBase.Save(user.ProxyRecord())
	if err != nil {
		return err
	}

	err = PocketBase.Save(user.LastAction().ProxyRecord())
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) DropItem(userId, itemId string) error {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return err
	}

	item, ok := user.Inventory().GetItemById(itemId)
	if !ok {
		return errors.New("item not found in inventory")
	}

	err = user.Inventory().DropItem(itemId)
	if err != nil {
		return err
	}

	if itemPrice := item.Price(); itemPrice > 0 {
		user.SetBalance(user.Balance() + itemPrice/2)

		err = PocketBase.Save(user.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) StartTimer(userId string) error {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return err
	}

	return user.Timer().Start()
}

func (g *Game) StopTimer(userId string) error {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return err
	}

	return user.Timer().Stop()
}

func (g *Game) GetTimeLeft(userId string) (int64, bool, types.DateTime, error) {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer().GetTimeLeft(), user.Timer().IsActive(), GameSettings.NextTimerResetDate(), nil
}

func (g *Game) GetAvailableActions(userId string) ([]ActionType, error) {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return nil, err
	}

	var actions []ActionType
	for t := range GameActions.AvailableActions(user) {
		actions = append(actions, t)
	}

	return actions, nil
}

func (g *Game) Context() context.Context {
	return g.ctx
}

func (g *Game) GetItemEffectVariants(userId, invItemId, effectId string) (any, error) {
	user, err := GameUsers.GetByID(userId)
	if err != nil {
		return nil, err
	}

	invItem, ok := user.Inventory().GetItemById(invItemId)
	if !ok {
		return nil, errors.New("inventory item not found")
	}

	return invItem.GetEffectVariants(effectId)
}
