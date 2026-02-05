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
	a.hookIds = make([]string, 2)

	a.hookIds[0] = ctx.App.OnRecordCreate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
	a.hookIds[1] = ctx.App.OnRecordUpdate(CollectionActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == a.ID() {
			e.Record.Set("custom_activity_filter", a.activityFilter)
		}
		return e.Next()
	})
}

func (a *LastUserActionRecord) Close(ctx AppContext) {
	ctx.App.OnRecordCreate(CollectionActions).Unbind(a.hookIds[0])
	ctx.App.OnRecordUpdate(CollectionActions).Unbind(a.hookIds[1])
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
