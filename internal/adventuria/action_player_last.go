package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type LastPlayerActionRecord struct {
	ActionRecordBase
	hookIds []string
}

func NewLastPlayerAction(ctx AppContext, playerId string) (*LastPlayerActionRecord, error) {
	a, err := getLastPlayerAction(ctx, playerId)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *LastPlayerActionRecord) Refetch(ctx AppContext) error {
	record, err := fetchLastPlayerAction(ctx, a.Player())
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

func getLastPlayerAction(ctx AppContext, playerId string) (*LastPlayerActionRecord, error) {
	record, err := fetchLastPlayerAction(ctx, playerId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	a := &LastPlayerActionRecord{}
	if record == nil {
		record = core.NewRecord(GameCollections.Get(schema.CollectionActions))
		record.Set(schema.ActionSchema.Type, ActionTypeNone)
		record.Set(schema.ActionSchema.CanMove, true)

		defaultWorld, ok := GameWorlds.GetDefault()
		if !ok {
			return nil, fmt.Errorf("getLastPlayerAction: default world not found")
		}

		firstCell, ok := GameCells.GetByOrder(defaultWorld.ID(), 0)
		if ok {
			record.Set(schema.ActionSchema.Cell, firstCell.ID())
		}
	}

	a.SetProxyRecord(record)
	a.SetPlayer(playerId)

	return a, nil
}

func fetchLastPlayerAction(ctx AppContext, playerId string) (*core.Record, error) {
	var record core.Record
	err := ctx.App.
		RecordQuery(schema.CollectionActions).
		Where(dbx.HashExp{schema.ActionSchema.Player: playerId}).
		AndWhere(dbx.NewExp("created > {:date}", dbx.Params{
			"date": GameSettings.CurrentSeasonDateStart(),
		})).
		OrderBy("created DESC", "rowid DESC").
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
