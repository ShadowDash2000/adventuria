package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type BaseGame struct {
	users       *cache.MemoryCache[string, User]
	pb          *pocketbase.PocketBase
	cells       *Cells
	items       *Items
	collections *collections.Collections
	settings    *Settings
}

func New() Game {
	game := &BaseGame{
		users: cache.NewMemoryCache[string, User](0, true),
	}

	game.pb = pocketbase.New()

	game.OnServe(func(se *core.ServeEvent) error {
		game.Init()
		return se.Next()
	})

	return game
}

func (g *BaseGame) OnServe(fn func(se *core.ServeEvent) error) {
	g.pb.OnServe().BindFunc(fn)
}

func (g *BaseGame) Start() error {
	return g.pb.Start()
}

func (g *BaseGame) PocketBase() *pocketbase.PocketBase {
	return g.pb
}

func (g *BaseGame) Cells() *Cells {
	return g.cells
}

func (g *BaseGame) Items() *Items {
	return g.items
}

func (g *BaseGame) Collections() *collections.Collections {
	return g.collections
}

func (g *BaseGame) Settings() *Settings {
	return g.settings
}

func (g *BaseGame) Init() {
	g.cells = NewCells(g)
	g.items = NewItems(g)
	g.collections = collections.NewCollections(g.pb)
	g.settings = NewSettings(g)
}

func (g *BaseGame) GetUser(userId string) (User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(g, userId)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *BaseGame) afterAction(user User) error {
	err := user.OnAfterAction().Trigger(&OnAfterActionEvent{})
	if err != nil {
		return err
	}

	err = user.LastAction().Save()
	if err != nil {
		return err
	}

	err = user.save()
	if err != nil {
		return err
	}

	/*	_, err = user.Inventory().ApplyEffectsByEvent(event)
		if err != nil {
			return err
		}*/

	return nil
}

func (g *BaseGame) GetNextStepType(userId string) (string, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	return user.GetNextStepType(), nil
}

func (g *BaseGame) DoAction(actionType, userId string, req ActionRequest) (*ActionResult, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	action, err := NewActionFromType(g, user, actionType)
	if err != nil {
		return nil, err
	}

	if ok := action.CanDo(); !ok {
		return nil, errors.New("cannot do action")
	}

	res, err := action.Do(req)
	if err != nil {
		return nil, err
	}

	err = user.OnAfterAction().Trigger(&OnAfterActionEvent{})
	if err != nil {
		return nil, err
	}

	err = user.LastAction().Save()
	if err != nil {
		return nil, err
	}

	err = user.save()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *BaseGame) UpdateAction(actionId string, comment string, userId string) error {
	actionsCollection, err := g.collections.Get(TableActions)
	if err != nil {
		return err
	}

	record := &core.Record{}
	err = g.pb.
		RecordQuery(actionsCollection).
		AndWhere(
			dbx.HashExp{
				"user": userId,
				"id":   actionId,
			},
		).
		AndWhere(
			dbx.Or(
				dbx.HashExp{"type": ActionTypeDone},
				dbx.HashExp{"type": ActionTypeDrop},
				dbx.HashExp{"type": ActionTypeReroll},
			),
		).
		Limit(1).
		One(record)
	if err != nil {
		return err
	}

	action := NewActionFromRecord(g, record)
	action.SetComment(comment)

	return action.Save()
}

func (g *BaseGame) GetLastAction(userId string) (bool, Action, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return false, nil, err
	}

	return user.CanDrop(), user.LastAction(), nil
}

func (g *BaseGame) UseItem(userId, itemId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	err = user.Inventory().UseItem(itemId)
	if err != nil {
		return err
	}

	onAfterItemUseEvent := &OnAfterItemUseEvent{
		ItemId: itemId,
	}
	err = user.OnAfterItemUse().Trigger(onAfterItemUseEvent)
	if err != nil {
		return err
	}

	err = g.afterAction(user)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) DropItem(userId, itemId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	err = user.Inventory().DropItem(itemId)
	if err != nil {
		return err
	}

	return nil
}

func (g *BaseGame) StartTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer().Start()
}

func (g *BaseGame) StopTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer().Stop()
}

func (g *BaseGame) GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer().GetTimeLeft(), user.Timer().IsActive(), g.settings.NextTimerResetDate(), nil
}
