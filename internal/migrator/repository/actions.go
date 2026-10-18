package repository

import (
	actionsRepo "adventuria/internal/adventuria/actions/repository"
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/migrator/model"
	"adventuria/pkg/pbtransaction"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type Actions struct {
	pb core.App
}

func NewActions(pb core.App) *Actions {
	return &Actions{pb: pb}
}

func (a *Actions) Create(ctx context.Context, actionCreate model.ActionImport) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, a.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionActions)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	err = actionsRepo.ActionToRecord(actionCreate.Action, record)
	if err != nil {
		return err
	}

	record.Id = actionCreate.ID
	record.SetRaw("created", actionCreate.Created)
	record.SetRaw("updated", actionCreate.Updated)

	return pb.SaveWithContext(ctx, record)
}
