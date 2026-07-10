package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
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

func (r *Repository) baseCellQuery(db dbx.Builder) *dbx.SelectQuery {
	return db.
		Select(
			fmt.Sprintf("%s.*", schema.CollectionCells),
			fmt.Sprintf(
				"ROW_NUMBER() OVER (PARTITION BY %s ORDER BY %s ASC, %s ASC) - 1 as local_order",
				schema.CellSchema.World,
				schema.CellSchema.Sort,
				schema.CellSchema.Id,
			),
			fmt.Sprintf(
				"ROW_NUMBER() OVER (ORDER BY %s ASC, %s ASC, %s ASC) - 1 as global_order",
				pbhelper.DotExpand(schema.CollectionWorlds, schema.WorldSchema.Sort),
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.Sort),
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.Id),
			),
		).
		From(schema.CollectionCells).
		LeftJoin(
			schema.CollectionWorlds,
			dbx.NewExp(pbhelper.Eq(
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.World),
				pbhelper.DotExpand(schema.CollectionWorlds, schema.WorldSchema.Id),
			)),
		).
		Where(dbx.HashExp{
			pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.Disabled): false,
		})
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := r.baseCellQuery(pb.DB()).Build()

	var record core.Record
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.HashExp{
			pbhelper.DotExpand("t", schema.CellSchema.Id): id,
		}).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return RecordToCell(&record), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	if len(ids) == 0 {
		return []*model.CellInfo{}, nil
	}

	subQuery := r.baseCellQuery(pb.DB()).Build()

	var records []*core.Record
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.In(
			pbhelper.DotExpand("t", schema.CellSchema.Id),
			pbhelper.SliceToAny(ids)...,
		)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToCells(records), nil
}

func (r *Repository) GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := r.baseCellQuery(pb.DB()).Build()

	var record core.Record
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.HashExp{"t." + schema.CellSchema.World: worldId}).
		OrderBy("t.local_order ASC").
		Limit(1).
		Offset(int64(order)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return RecordToCell(&record), nil
}

func (r *Repository) GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := r.baseCellQuery(pb.DB()).Build()

	var record core.Record
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		OrderBy("t.global_order ASC").
		Limit(1).
		Offset(int64(order)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return RecordToCell(&record), nil
}

func (r *Repository) GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := r.baseCellQuery(pb.DB()).Build()

	var records []*core.Record
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		OrderBy("t.global_order ASC").
		Where(dbx.HashExp{
			pbhelper.DotExpand("t", schema.CellSchema.Type): string(t),
		}).
		Bind(subQuery.Params()).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToCells(records), nil
}

func (r *Repository) CountLocal(ctx context.Context, worldId string) (int, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var count int
	err := pb.DB().
		Select("count(*)").
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.World:    worldId,
			schema.CellSchema.Disabled: false,
		}).
		WithContext(ctx).
		Row(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) CountGlobal(ctx context.Context) (int, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var count int
	err := pb.DB().
		Select("count(*)").
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.Disabled: false,
		}).
		WithContext(ctx).
		Row(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
