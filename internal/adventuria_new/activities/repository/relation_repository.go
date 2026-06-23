package repository

import (
	"adventuria/pkg/pbtransaction"
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type RelationRepository struct {
	pb core.App
}

func NewRelationRepository(pb core.App) *RelationRepository {
	return &RelationRepository{pb: pb}
}

func (r *RelationRepository) SyncRelations(
	ctx context.Context, collection, activityField, relationField, activityId string, relationIds []string,
) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(collection).
		WithContext(ctx).
		Where(dbx.HashExp{activityField: activityId}).
		All(&records)
	if err != nil {
		return err
	}

	toDelete := make(map[string]*core.Record, len(records))
	for _, record := range records {
		toDelete[record.GetString(relationField)] = record
	}

	for _, relationId := range relationIds {
		if _, ok := toDelete[relationId]; ok {
			delete(toDelete, relationId)
		} else {
			err = r.createRelation(ctx, collection, activityField, relationField, activityId, relationId)
			if err != nil {
				return err
			}
		}
	}

	for _, record := range toDelete {
		err = pb.DeleteWithContext(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RelationRepository) DeleteAllRelations(
	ctx context.Context, collection, activityField, activityId string,
) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(collection).
		WithContext(ctx).
		Where(dbx.HashExp{activityField: activityId}).
		All(&records)
	if err != nil {
		return err
	}

	for _, record := range records {
		err = pb.DeleteWithContext(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RelationRepository) createRelation(
	ctx context.Context, collection, activityField, relationField, activityId, relationId string,
) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	coll, err := pb.FindCollectionByNameOrId(collection)
	if err != nil {
		return err
	}

	record := core.NewRecord(coll)
	record.Set(activityField, activityId)
	record.Set(relationField, relationId)

	return pb.SaveWithContext(ctx, record)
}
