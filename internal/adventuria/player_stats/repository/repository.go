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

func (r *Repository) Create(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayerStats)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerStatsToRecord(stats, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerStats(record), nil
}

func (r *Repository) Update(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlayerStats, stats.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerStatsNotFound
		}
		return nil, err
	}

	PlayerStatsToRecord(stats, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerStats(record), nil
}

func (r *Repository) Save(ctx context.Context, stats *model.PlayerStats) (*model.PlayerStats, error) {
	if stats.IsNew() {
		return r.Create(ctx, stats)
	}

	return r.Update(ctx, stats)
}

func (r *Repository) GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerStats, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayerStats).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.PlayerProgressSchema.Player: playerId,
			schema.PlayerStatsSchema.Season:    seasonId,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerStatsNotFound
		}
		return nil, err
	}

	return RecordToPlayerStats(&record), nil
}
