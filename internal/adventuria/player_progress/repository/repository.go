package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
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

func (r *Repository) Create(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayersProgress)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerProgressToRecord(progress, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerProgress(record), nil
}

func (r *Repository) Update(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlayersProgress, progress.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProgressNotFound
		}
		return nil, err
	}

	PlayerProgressToRecord(progress, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerProgress(record), nil
}

func (r *Repository) Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	if progress.IsNew() {
		return r.Create(ctx, progress)
	}

	return r.Update(ctx, progress)
}

func (r *Repository) GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayersProgress).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.PlayerProgressSchema.Player: playerId},
			dbx.HashExp{schema.PlayerProgressSchema.Season: seasonId},
		)).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProgressNotFound
		}
		return nil, err
	}

	progress := RecordToPlayerProgress(&record)

	return progress, nil
}

func (r *Repository) ChangeBalance(ctx context.Context, id string, amount int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	_, err := pb.DB().
		NewQuery(fmt.Sprintf(
			"UPDATE %s SET %s = %s + {:amount} WHERE %s = {:id}",
			schema.CollectionPlayersProgress,
			schema.PlayerProgressSchema.Balance,
			schema.PlayerProgressSchema.Balance,
			schema.PlayerProgressSchema.Id,
		)).
		Bind(dbx.Params{
			"amount": amount,
			"id":     id,
		}).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) NotifyChange(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayersProgress).
		WithContext(ctx).
		Where(dbx.HashExp{schema.PlayerProgressSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrProgressNotFound
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
