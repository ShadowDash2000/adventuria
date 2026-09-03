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

func (r *Repository) GetByID(ctx context.Context, id string) (*model.PlayerInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerNotFound
		}
		return nil, err
	}

	return RecordToPlayerInfo(&record), nil
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

func (r *Repository) IsDisabled(ctx context.Context, id string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var isDisabled bool
	err := pb.DB().
		Select(schema.PlayerSchema.Disabled).
		From(schema.CollectionPlayers).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		WithContext(ctx).
		Row(&isDisabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.ErrPlayerNotFound
		}
		return false, err
	}

	return isDisabled, nil
}

func (r *Repository) IsDebugEnabled(ctx context.Context, id string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var isDebugEnabled bool
	err := pb.DB().
		Select(schema.PlayerSchema.Debug).
		From(schema.CollectionPlayers).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		WithContext(ctx).
		Row(&isDebugEnabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.ErrPlayerNotFound
		}
		return false, err
	}

	return isDebugEnabled, nil
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
