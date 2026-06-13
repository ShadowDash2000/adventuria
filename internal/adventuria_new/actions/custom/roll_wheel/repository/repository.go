package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/activities/repository"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
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

func (r *Repository) GetDetailedByIDs(ctx context.Context, ids []string) ([]*model.ActivityViewDetailed, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.In(
			schema.ActivitySchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	errs := pb.ExpandRecords(records, []string{
		schema.ActivitySchema.Platforms,
		schema.ActivitySchema.Developers,
		schema.ActivitySchema.Publishers,
		schema.ActivitySchema.Genres,
		schema.ActivitySchema.Tags,
		schema.ActivitySchema.Themes,
	}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
	}

	activities := make([]*model.ActivityViewDetailed, len(records))
	for i, record := range records {
		activities[i] = repository.RecordToActivityViewDetailed(record)
	}

	return activities, nil
}
