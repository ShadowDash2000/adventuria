package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) UpdateActiveCellAndNextShiftByID(ctx context.Context, id, cellId string, nextShiftChangeAt time.Time) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionCellEventsSchedule,
			dbx.Params{
				schema.CellEventsScheduleSchema.ActiveCell:        cellId,
				schema.CellEventsScheduleSchema.NextShiftChangeAt: nextShiftChangeAt,
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

func (r *Repository) NotifyChange(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionCellEventsSchedule).
		WithContext(ctx).
		Where(dbx.HashExp{schema.CellEventsScheduleSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrCellEventScheduleNotFound
		}
		return err
	}

	event := &core.ModelEvent{
		App:     pb,
		Context: ctx,
		Type:    core.ModelEventTypeUpdate,
	}
	event.Model = &record

	return pb.OnModelAfterUpdateSuccess().Trigger(event)
}
