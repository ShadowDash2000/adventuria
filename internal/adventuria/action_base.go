package adventuria

import (
	"adventuria/pkg/cache"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// ensure Action implements Action
var _ Action = (*ActionBase)(nil)

type ActionBase struct {
	t    ActionType
	user User
}

func NewActionFromType(actionType ActionType) (Action, error) {
	actionCreator, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionCreator()

	return action, nil
}

func (a *ActionBase) Type() ActionType {
	return a.t
}

func (a *ActionBase) CanDo(_ User) bool {
	panic("implement me")
}

func (a *ActionBase) Do(_ User, _ ActionRequest) (*ActionResult, error) {
	panic("implement me")
}

func (a *ActionBase) setType(t ActionType) {
	a.t = t
}

type ActionRecordBase struct {
	core.BaseRecordProxy
	gameFilter CustomGameFilter
}

func NewActionRecord() ActionRecord {
	actionRecord := &ActionRecordBase{}
	actionRecord.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))

	return actionRecord
}

func NewActionRecordFromRecord(record *core.Record) ActionRecord {
	a := &ActionRecordBase{}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionRecordBase) SetProxyRecord(record *core.Record) {
	a.BaseRecordProxy.SetProxyRecord(record)
	a.UnmarshalJSONField("game_filter", &a.gameFilter)
}

func (a *ActionRecordBase) ID() string {
	return a.Id
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

func (a *ActionRecordBase) setCell(cellId string) {
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

func (a *ActionRecordBase) CustomGameFilter() *CustomGameFilter {
	return &a.gameFilter
}

var _ cache.Closable = (*LastUserActionRecord)(nil)

type LastUserActionRecord struct {
	ActionRecordBase
	hookIds []string
}

func NewLastUserAction(userId string) (*LastUserActionRecord, error) {
	a, err := getLastUserAction(userId)
	if err != nil {
		return nil, err
	}

	a.bindHooks()

	return a, nil
}

func (a *LastUserActionRecord) bindHooks() {
	a.hookIds = make([]string, 3)

	a.hookIds[0] = PocketBase.OnRecordAfterCreateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != a.User() {
			return e.Next()
		}

		a.SetProxyRecord(e.Record)

		return e.Next()
	})
	a.hookIds[1] = PocketBase.OnRecordAfterUpdateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id != a.ID() {
			return e.Next()
		}

		a.SetProxyRecord(e.Record)

		return e.Next()
	})
	a.hookIds[2] = PocketBase.OnRecordAfterDeleteSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != a.User() {
			return e.Next()
		}

		record, err := fetchLastUserAction(a.User())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				a.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
				a.SetType(ActionTypeNone)
				a.SetCanMove(true)

				return e.Next()
			} else {
				return err
			}
		}

		a.SetProxyRecord(record)

		return e.Next()
	})
}

func (a *LastUserActionRecord) Close() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionActions).Unbind(a.hookIds[0])
	PocketBase.OnRecordAfterUpdateSuccess(CollectionActions).Unbind(a.hookIds[1])
	PocketBase.OnRecordAfterDeleteSuccess(CollectionActions).Unbind(a.hookIds[2])
}

func getLastUserAction(userId string) (*LastUserActionRecord, error) {
	record, err := fetchLastUserAction(userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	a := &LastUserActionRecord{}
	if errors.Is(err, sql.ErrNoRows) {
		a.SetProxyRecord(core.NewRecord(GameCollections.Get(CollectionActions)))
		a.SetType(ActionTypeNone)
		a.SetCanMove(true)
		firstCell, ok := GameCells.GetByOrder(0)
		if ok {
			a.setCell(firstCell.ID())
		}
	} else {
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
