package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/collections"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GameOpts struct {
	users cache.Cache[string, User]
}

type GameOptsFunc func(*GameOpts)

func WithUsersCache(users cache.Cache[string, User]) GameOptsFunc {
	return func(opts *GameOpts) {
		opts.users = users
	}
}

var (
	PocketBase      core.App
	GameCells       *Cells
	GameItems       *Items
	GameCollections *collections.Collections
	GameSettings    *Settings
)

type Game struct {
	opts GameOpts
	pb   *pocketbase.PocketBase
	ef   *EffectVerifier
}

func New(opts ...GameOptsFunc) *Game {
	game := &Game{
		opts: GameOpts{},
	}

	for _, opt := range opts {
		opt(&game.opts)
	}

	if game.opts.users == nil {
		game.opts.users = cache.NewMemoryCache[string, User](time.Hour, false)
	}

	return game
}

func (g *Game) Start(fn func(se *core.ServeEvent) error) error {
	g.pb = pocketbase.New()
	PocketBase = g.pb

	migratecmd.MustRegister(g.pb, g.pb.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	g.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if err := g.init(); err != nil {
			return err
		}
		return se.Next()
	})

	g.pb.OnServe().BindFunc(fn)

	return g.pb.Start()
}

func (g *Game) init() error {
	var err error

	GameCollections = collections.NewCollections(PocketBase)
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
	return nil
}

func (g *Game) GetUser(userId string) (User, error) {
	user, ok := g.opts.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId)
	if err != nil {
		return nil, err
	}

	g.opts.users.Set(userId, user)
	return user, nil
}

func (g *Game) GetUserByName(name string) (User, error) {
	for _, user := range g.opts.users.GetAll() {
		if name == user.Name() {
			return user, nil
		}
	}

	user, err := NewUserFromName(name)
	if err != nil {
		return nil, err
	}

	g.opts.users.Set(user.ID(), user)
	return user, nil
}

func (g *Game) DoAction(actionType ActionType, userId string, req ActionRequest) (*ActionResult, error) {
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
		PocketBase.Logger().Error("Failed to complete user action", "error", err)
		return nil, err
	} else if res.Error != "" {
		return nil, errors.New(res.Error)
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

// UpdateAction TODO: reimplement, method expired
func (g *Game) UpdateAction(actionId string, comment string, userId string) error {
	record := &core.Record{}
	err := PocketBase.
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

func (g *Game) UseItem(userId, itemId string) error {
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

func (g *Game) DropItem(userId, itemId string) error {
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

func (g *Game) StartTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer().Start()
}

func (g *Game) StopTimer(userId string) error {
	user, err := g.GetUser(userId)
	if err != nil {
		return err
	}

	return user.Timer().Stop()
}

func (g *Game) GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return 0, false, types.DateTime{}, err
	}

	return user.Timer().GetTimeLeft(), user.Timer().IsActive(), GameSettings.NextTimerResetDate(), nil
}

func (g *Game) GetAvailableActions(userId string) ([]ActionType, error) {
	user, err := g.GetUser(userId)
	if err != nil {
		return nil, err
	}

	var actions []ActionType
	for t, _ := range actionsList {
		action, _ := NewActionFromType(user, t)

		if action.CanDo() {
			actions = append(actions, action.Type())
		}
	}

	return actions, nil
}
