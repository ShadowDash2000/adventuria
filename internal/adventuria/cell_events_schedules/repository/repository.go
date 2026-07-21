package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) UpdateActiveCellByID(ctx context.Context, id, cellId string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionCellEventsSchedule,
			dbx.Params{
				schema.CellEventsScheduleSchema.ActiveCell:      cellId,
				schema.CellEventsScheduleSchema.LastShiftChange: types.NowDateTime(),
			},
			dbx.HashExp{
				schema.CellEventsScheduleSchema.Id: id,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrCellEventScheduleNotFound
	}

	return nil
}

func (r *Repository) GetByActiveCellID(ctx context.Context, activeCellId string) (*model.CellEventSchedule, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionCellEventsSchedule).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.CellEventsScheduleSchema.ActiveCell: activeCellId,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellEventScheduleNotFound
		}
		return nil, err
	}

	return RecordToCellEventSchedule(&record), nil
}

func (r *Repository) GetIDByActiveCellID(ctx context.Context, activeCellId string) (string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var id string
	err := pb.DB().
		Select(schema.CellEventsScheduleSchema.Id).
		From(schema.CollectionCellEventsSchedule).
		Where(dbx.HashExp{
			schema.CellEventsScheduleSchema.ActiveCell: activeCellId,
		}).
		Limit(1).
		WithContext(ctx).
		Row(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrCellEventScheduleNotFound
		}
		return "", err
	}

	return id, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.CellEventSchedule, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionCellEventsSchedule).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToCellEventSchedules(records), nil
}
