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
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Select(schema.PlayerSchema.Id).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repository) GetAllIDs(ctx context.Context) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Select(schema.PlayerSchema.Id).
		All(&records)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(records))
	for i, record := range records {
		ids[i] = record.Id
	}

	return ids, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.PlayerInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToPlayerInfos(records), nil
}

func (r *Repository) NotifyChange(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrPlayerNotFound
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
