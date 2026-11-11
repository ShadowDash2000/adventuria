package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, task *model.PlayerTaskInfo) (*model.PlayerTaskInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayerTasks)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerTaskToRecord(task, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerTask(record), nil
}

func (r *Repository) Update(ctx context.Context, task *model.PlayerTaskInfo) (*model.PlayerTaskInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlayerTasks, task.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerTaskNotFound
		}
		return nil, err
	}

	PlayerTaskToRecord(task, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerTask(record), nil
}
