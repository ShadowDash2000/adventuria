package adventuria

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ActionBase struct {
	t    ActionType
	user User
}

func NewActionFromType(user User, actionType ActionType) (Action, error) {
	actionCreator, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionCreator()
	action.setUser(user)

	return action, nil
}

func (a *ActionBase) User() User {
	return a.user
}

func (a *ActionBase) Type() ActionType {
	return a.t
}

func (a *ActionBase) CanDo() bool {
	panic("implement me")
}

func (a *ActionBase) Do(_ ActionRequest) (*ActionResult, error) {
	panic("implement me")
}

func (a *ActionBase) setType(t ActionType) {
	a.t = t
}

func (a *ActionBase) setUser(user User) {
	a.user = user
}

type ActionRecordBase struct {
	core.BaseRecordProxy
}

func NewActionRecordFromType(actionType ActionType) (ActionRecord, error) {
	_, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	actionRecord := &ActionRecordBase{}
	actionRecord.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
	actionRecord.SetType(actionType)

	return actionRecord, nil
}

// TODO: remove unused
func NewActionFromRecord(record *core.Record) ActionRecord {
	a := &ActionRecordBase{}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionRecordBase) ID() string {
	return a.Id
}

func (a *ActionRecordBase) Save() error {
	return PocketBase.Save(a)
}

func (a *ActionRecordBase) User() string {
	return a.GetString("user")
}

func (a *ActionRecordBase) SetUser(id string) {
	a.Set("user", id)
}

func (a *ActionRecordBase) CellId() string {
	return a.GetString("cell")
}

func (a *ActionRecordBase) SetCell(cellId string) {
	a.Set("cell", cellId)
}

func (a *ActionRecordBase) Comment() string {
	return a.GetString("comment")
}

func (a *ActionRecordBase) SetComment(comment string) {
	a.Set("comment", comment)
}

func (a *ActionRecordBase) Game() string {
	return a.GetString("game")
}

func (a *ActionRecordBase) SetGame(id string) {
	a.Set("game", id)
}

func (a *ActionRecordBase) Type() ActionType {
	return ActionType(a.GetString("type"))
}

func (a *ActionRecordBase) SetType(t ActionType) {
	a.Set("type", string(t))
}

func (a *ActionRecordBase) SetNotAffectNextStep(b bool) {
	a.Set("notAffectNextStep", b)
}

func (a *ActionRecordBase) DiceRoll() int {
	return a.GetInt("diceRoll")
}

func (a *ActionRecordBase) SetDiceRoll(roll int) {
	a.Set("diceRoll", roll)
}

func (a *ActionRecordBase) ItemsUsed() []string {
	return a.GetStringSlice("itemsUsed")
}

func (a *ActionRecordBase) SetItemsUsed(items []string) {
	a.Set("itemsUsed", items)
}

func (a *ActionRecordBase) ItemsList() ([]string, error) {
	var items []string
	return items, a.UnmarshalJSONField("items_list", &items)
}

func (a *ActionRecordBase) SetItemsList(items []string) {
	a.Set("items_list", items)
}

func (a *ActionRecordBase) CanMove() bool {
	return a.GetBool("can_move")
}

func (a *ActionRecordBase) SetCanMove(b bool) {
	a.Set("can_move", b)
}

func NewLastUserAction(userId string) (ActionRecord, error) {
	a, err := getLastUserAction(userId)
	if err != nil {
		return nil, err
	}
	actionBindHooks(a)

	return a, nil
}

func actionBindHooks(action ActionRecord) {
	PocketBase.OnRecordAfterCreateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != action.User() {
			return e.Next()
		}

		actionType := ActionType(e.Record.GetString("type"))
		if actionType == action.Type() {
			action.SetProxyRecord(e.Record)
			return e.Next()
		}

		a, err := NewActionRecordFromType(actionType)
		if err != nil {
			return err
		}

		action = a
		action.SetProxyRecord(e.Record)

		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == action.ID() {
			return e.Next()
		}

		actionType := ActionType(e.Record.GetString("type"))
		if actionType == action.Type() {
			action.SetProxyRecord(e.Record)
			return e.Next()
		}

		a, err := NewActionRecordFromType(actionType)
		if err != nil {
			return err
		}

		action = a
		action.SetProxyRecord(e.Record)

		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != action.User() {
			return e.Next()
		}

		a, err := getLastUserAction(action.User())
		if err != nil {
			return err
		}

		action = a

		return e.Next()
	})
}

func getLastUserAction(userId string) (ActionRecord, error) {
	record, err := fetchLastUserAction(userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var a ActionRecord
	if errors.Is(err, sql.ErrNoRows) {
		a, err = NewActionRecordFromType(ActionTypeNone)
		if err != nil {
			return nil, err
		}
		a.SetCanMove(true)
	} else {
		a, err = NewActionRecordFromType(ActionType(record.GetString("type")))
		if err != nil {
			return nil, err
		}

		a.SetProxyRecord(record)
	}

	a.SetUser(userId)

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
