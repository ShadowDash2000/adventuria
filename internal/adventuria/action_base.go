package adventuria

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ActionBase struct {
	core.BaseRecordProxy
	user User
}

func NewActionFromType(user User, actionType ActionType) (Action, error) {
	actionCreator, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionCreator()

	action.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
	action.setUser(user)
	action.SetType(actionType)

	return action, nil
}

func NewActionFromRecord(record *core.Record) Action {
	a := &ActionBase{}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionBase) CanDo() bool {
	panic("implement me")
}

func (a *ActionBase) NextAction() ActionType {
	panic("implement me")
}

func (a *ActionBase) Do(_ ActionRequest) (*ActionResult, error) {
	panic("implement me")
}

func (a *ActionBase) setUser(user User) {
	a.user = user
	a.Set("user", user.ID())
}

func (a *ActionBase) ID() string {
	return a.Id
}

func (a *ActionBase) Save() error {
	return PocketBase.Save(a)
}

func (a *ActionBase) User() User {
	return a.user
}

func (a *ActionBase) UserId() string {
	return a.GetString("user")
}

func (a *ActionBase) CellId() string {
	return a.GetString("cell")
}

func (a *ActionBase) SetCell(cellId string) {
	a.Set("cell", cellId)
}

func (a *ActionBase) Comment() string {
	return a.GetString("comment")
}

func (a *ActionBase) SetComment(comment string) {
	a.Set("comment", comment)
}

func (a *ActionBase) Game() string {
	return a.GetString("game")
}

func (a *ActionBase) SetGame(id string) {
	a.Set("game", id)
}

func (a *ActionBase) Type() ActionType {
	return ActionType(a.GetString("type"))
}

func (a *ActionBase) SetType(t ActionType) {
	a.Set("type", string(t))
}

func (a *ActionBase) SetNotAffectNextStep(b bool) {
	a.Set("notAffectNextStep", b)
}

func (a *ActionBase) DiceRoll() int {
	return a.GetInt("diceRoll")
}

func (a *ActionBase) SetDiceRoll(roll int) {
	a.Set("diceRoll", roll)
}

func (a *ActionBase) ItemsUsed() []string {
	return a.GetStringSlice("itemsUsed")
}

func (a *ActionBase) SetItemsUsed(items []string) {
	a.Set("itemsUsed", items)
}

func (a *ActionBase) Seed() int {
	return a.GetInt("seed")
}

func (a *ActionBase) SetSeed(seed int) {
	a.Set("seed", seed)
}

type UserAction struct {
	ActionBase
}

func NewLastUserAction(user User) (Action, error) {
	a, err := getLastUserAction(user)
	if err != nil {
		return nil, err
	}
	actionBindHooks(a)

	return a, nil
}

func actionBindHooks(action Action) {
	PocketBase.OnRecordAfterCreateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == action.UserId() {
			action.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == action.ID() {
			action.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == action.UserId() {
			a, err := getLastUserAction(action.User())
			if err != nil {
				return err
			}

			action = a
		}
		return e.Next()
	})
}

func getLastUserAction(user User) (Action, error) {
	record, err := fetchLastUserAction(user.ID())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var a Action
	if errors.Is(err, sql.ErrNoRows) {
		a, err = NewActionFromType(user, "none")
		if err != nil {
			return nil, err
		}
	} else {
		a, err = NewActionFromType(user, ActionType(record.GetString("type")))
		if err != nil {
			return nil, err
		}

		a.SetProxyRecord(record)
	}

	return a, nil
}

func fetchLastUserAction(userId string) (*core.Record, error) {
	actions, err := PocketBase.FindRecordsByFilter(
		CollectionActions,
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": userId},
	)
	if err != nil {
		return nil, err
	}

	if len(actions) == 0 {
		return nil, sql.ErrNoRows
	}

	return actions[0], nil
}
