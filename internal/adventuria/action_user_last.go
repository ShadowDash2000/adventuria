package adventuria

import (
	"adventuria/pkg/cache"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

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
	a.hookIds = make([]string, 5)

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
			}

			return err
		}

		a.SetProxyRecord(record)

		return e.Next()
	})
	a.hookIds[3] = PocketBase.OnRecordCreate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
	a.hookIds[4] = PocketBase.OnRecordUpdate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
}

func (a *LastUserActionRecord) Close() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionActions).Unbind(a.hookIds[0])
	PocketBase.OnRecordAfterUpdateSuccess(CollectionActions).Unbind(a.hookIds[1])
	PocketBase.OnRecordAfterDeleteSuccess(CollectionActions).Unbind(a.hookIds[2])
	PocketBase.OnRecordCreate(CollectionActions).Unbind(a.hookIds[3])
	PocketBase.OnRecordUpdate(CollectionActions).Unbind(a.hookIds[4])
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
