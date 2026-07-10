package repository

import (
	igdbparser "adventuria/internal/adventuria/games/igdb"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) TableReferenceToID(ctx context.Context, reference igdbparser.TableReferenceSingle) (string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var id string
	err := pb.RecordQuery(reference.TableName).
		WithContext(ctx).
		Select(reference.PrimaryKey).
		Where(dbx.HashExp{
			reference.SearchKey: reference.Id,
		}).
		Row(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) TableReferenceToIDs(ctx context.Context, reference igdbparser.TableReference) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var ids []string
	err := pb.RecordQuery(reference.TableName).
		WithContext(ctx).
		Select(reference.PrimaryKey).
		Where(dbx.In(
			reference.SearchKey,
			pbhelper.SliceToAny(reference.Ids)...,
		)).
		Column(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
