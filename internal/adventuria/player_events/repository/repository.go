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

func (r *Repository) Create(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayerEvents)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerEventToRecord(playerEvent, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerEvent(record), nil
}

func (r *Repository) Update(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayerEvents).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.PlayerEventsSchema.Id: playerEvent.ID(),
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerEventNotFound
		}
		return nil, err
	}

	PlayerEventToRecord(playerEvent, &record)
	err = pb.SaveWithContext(ctx, &record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerEvent(&record), nil
}
