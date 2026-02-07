package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type LastUserActionRecord struct {
	ActionRecordBase
	hookIds []string
}

func NewLastUserAction(ctx AppContext, userId string) (*LastUserActionRecord, error) {
	a, err := getLastUserAction(ctx, userId)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *LastUserActionRecord) Refetch(ctx AppContext) error {
	record, err := fetchLastUserAction(ctx, a.User())
	if err != nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		a.SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionActions)))
		a.SetType(ActionTypeNone)
		a.SetCanMove(true)
	} else {
		a.SetProxyRecord(record)
	}

	return nil
}

func getLastUserAction(ctx AppContext, userId string) (*LastUserActionRecord, error) {
	record, err := fetchLastUserAction(ctx, userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	a := &LastUserActionRecord{}
	if errors.Is(err, sql.ErrNoRows) {
		a.SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionActions)))
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
	var record core.Record
	err := ctx.App.
		RecordQuery(schema.CollectionActions).
		Where(dbx.HashExp{schema.ActionSchema.User: userId}).
		OrderBy("created DESC", "rowid DESC").
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
