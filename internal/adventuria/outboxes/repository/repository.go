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

func (r *Repository) Create(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionsOutbox)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	OutboxToRecord(outbox, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToOutbox(record), nil
}

func (r *Repository) Update(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionsOutbox, outbox.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrOutboxNotFound
		}
		return nil, err
	}

	OutboxToRecord(outbox, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToOutbox(record), nil
}

func (r *Repository) Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error) {
	if outbox.IsNew() {
		return r.Create(ctx, outbox)
	}

	return r.Update(ctx, outbox)
}

func (r *Repository) GetAndLockNextPending(ctx context.Context) (*model.OutboxInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pbtransaction.RunInTransaction(ctx, pb, func(ctx context.Context, txApp core.App) error {
		res, err := txApp.DB().
			Update(
				schema.CollectionsOutbox,
				dbx.Params{
					schema.OutboxSchema.Status: model.OutboxStatusProcessing,
				},
				dbx.NewExp(
					fmt.Sprintf(
						"%s = (SELECT %s FROM %s WHERE %s = {:status} ORDER BY created ASC LIMIT 1)",
						schema.OutboxSchema.Id,
						schema.OutboxSchema.Id,
						schema.CollectionsOutbox,
						schema.OutboxSchema.Status,
					),
					dbx.Params{
						"status": model.OutboxStatusPending,
					},
				),
			).
			WithContext(ctx).
			Execute()
		if err != nil {
			return err
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return errs.ErrNoPendingOutbox
		}

		return txApp.RecordQuery(schema.CollectionsOutbox).
			WithContext(ctx).
			Where(dbx.HashExp{
				schema.OutboxSchema.Status: model.OutboxStatusProcessing,
			}).
			OrderBy("updated DESC").
			One(&record)
	})
	if err != nil {
		return nil, err
	}

	return RecordToOutbox(&record), nil
}
