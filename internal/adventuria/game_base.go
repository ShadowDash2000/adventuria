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
	users *cache.MemoryCache[string, User]
	pb    *pocketbase.PocketBase
}

func New() Game {
	game := &BaseGame{
		users: cache.NewMemoryCache[string, User](0, true),
	}

	game.pb = pocketbase.New()
	PocketBase = game.pb

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

func (g *BaseGame) Init() {
	GameCells = NewCells()
	GameItems = NewItems()
	GameCollections = collections.NewCollections(PocketBase)
	GameSettings = NewSettings()
}

func (g *BaseGame) GetUser(userId string) (User, error) {
	user, ok := g.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId)
	if err != nil {
		return nil, err
	}

	g.users.Set(userId, user)
	return user, nil
}

func (g *BaseGame) GetUserByName(name string) (User, error) {
	for _, user := range g.users.GetAll() {
		if name == user.Name() {
			return user, nil
		}
	}

	user, err := NewUserFromName(name)
	if err != nil {
		return nil, err
	}

	g.users.Set(user.ID(), user)
	return user, nil
}

func (g *BaseGame) NextActionType(userId string) (ActionType, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return "", err
	}

	return user.NextAction(), nil
}

func (g *BaseGame) DoAction(actionType ActionType, userId string, req ActionRequest) (*ActionResult, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	action, err := NewActionFromType(user, actionType)
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
	record := &core.Record{}
	err := g.pb.
		RecordQuery(GameCollections.Get(CollectionActions)).
		AndWhere(
			dbx.HashExp{
				"user": userId,
				"id":   actionId,
			},
		).
		AndWhere(
			dbx.Or(
				// TODO get rid of hard coded types
				dbx.HashExp{"type": "done"},
				dbx.HashExp{"type": "drop"},
				dbx.HashExp{"type": "reroll"},
			),
		).
		Limit(1).
		One(record)
	if err != nil {
		return err
	}

	action := NewActionFromRecord(record)
	action.SetComment(comment)

	return action.Save()
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

	err = user.OnAfterAction().Trigger(&OnAfterActionEvent{})
	if err != nil {
		return err
	}

	err = user.save()
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

	return user.Timer().GetTimeLeft(), user.Timer().IsActive(), GameSettings.NextTimerResetDate(), nil
}
