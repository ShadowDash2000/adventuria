package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionActionEvents).
		WithContext(ctx).
		Select(fmt.Sprintf("%s.*", schema.CollectionActionEvents)).
		Where(dbx.HashExp{
			pbhelper.DotExpand(
				schema.CollectionCellEventsSchedule,
				schema.CellEventsScheduleSchema.ActiveCell,
			): activeCellId,
		}).
		InnerJoin(
			schema.CollectionCellEventsSchedule,
			dbx.NewExp(pbhelper.Eq(
				pbhelper.DotExpand(schema.CollectionActionEvents, schema.ActionEventsSchema.Id),
				pbhelper.DotExpand(schema.CollectionCellEventsSchedule, schema.CellEventsScheduleSchema.ActionEvent),
			)),
		).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActionEventNotFound
		}
		return nil, err
	}

	return RecordToActionEvent(&record), err
}
