package adventuria

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

var _ Closable = (*LastUserActionRecord)(nil)

type LastUserActionRecord struct {
	ActionRecordBase
	hookIds []string
}

func NewLastUserAction(ctx AppContext, userId string) (*LastUserActionRecord, error) {
	a, err := getLastUserAction(ctx, userId)
	if err != nil {
		return nil, err
	}

	a.bindHooks(ctx)

	return a, nil
}

func (a *LastUserActionRecord) bindHooks(ctx AppContext) {
	a.hookIds = make([]string, 5)

	a.hookIds[0] = ctx.App.OnRecordAfterCreateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != a.User() {
			return e.Next()
		}

		a.SetProxyRecord(e.Record)

		return e.Next()
	})
	a.hookIds[1] = ctx.App.OnRecordAfterUpdateSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id != a.ID() {
			return e.Next()
		}

		a.SetProxyRecord(e.Record)

		return e.Next()
	})
	a.hookIds[2] = ctx.App.OnRecordAfterDeleteSuccess(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId != a.User() {
			return e.Next()
		}

		record, err := fetchLastUserAction(AppContext{App: e.App}, a.User())
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
	a.hookIds[3] = ctx.App.OnRecordCreate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
	a.hookIds[4] = ctx.App.OnRecordUpdate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
}

func (a *LastUserActionRecord) Close(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(CollectionActions).Unbind(a.hookIds[0])
	ctx.App.OnRecordAfterUpdateSuccess(CollectionActions).Unbind(a.hookIds[1])
	ctx.App.OnRecordAfterDeleteSuccess(CollectionActions).Unbind(a.hookIds[2])
	ctx.App.OnRecordCreate(CollectionActions).Unbind(a.hookIds[3])
	ctx.App.OnRecordUpdate(CollectionActions).Unbind(a.hookIds[4])
}

func getLastUserAction(ctx AppContext, userId string) (*LastUserActionRecord, error) {
	record, err := fetchLastUserAction(ctx, userId)
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

func fetchLastUserAction(ctx AppContext, userId string) (*core.Record, error) {
	actions, err := ctx.App.FindRecordsByFilter(
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
