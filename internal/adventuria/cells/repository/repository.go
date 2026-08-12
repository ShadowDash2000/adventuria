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

func baseCellQuery(db dbx.Builder) *dbx.SelectQuery {
	return db.
		Select(
			fmt.Sprintf("%s.*", schema.CollectionCells),
			fmt.Sprintf(
				"ROW_NUMBER() OVER (PARTITION BY %s ORDER BY %s ASC, %s ASC) - 1 as local_order",
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.World),
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.Sort),
				pbhelper.DotExpand(schema.CollectionCells, schema.CellSchema.Id),
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

	subQuery := baseCellQuery(pb.DB()).Build()

	var cell cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.HashExp{
			pbhelper.DotExpand("t", schema.CellSchema.Id): id,
		}).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&cell)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return cellRowToCellInfo(&cell), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	if len(ids) == 0 {
		return []*model.CellInfo{}, nil
	}

	subQuery := baseCellQuery(pb.DB()).Build()

	var cells []cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.In(
			pbhelper.DotExpand("t", schema.CellSchema.Id),
			pbhelper.SliceToAny(ids)...,
		)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		All(&cells)
	if err != nil {
		return nil, err
	}

	return cellRowsToCellInfos(cells), nil
}

func (r *Repository) GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := baseCellQuery(pb.DB()).Build()

	var cell cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		Where(dbx.HashExp{"t." + schema.CellSchema.World: worldId}).
		OrderBy("t.local_order ASC").
		Limit(1).
		Offset(int64(order)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&cell)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return cellRowToCellInfo(&cell), nil
}

func (r *Repository) GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := baseCellQuery(pb.DB()).Build()

	var cell cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		OrderBy("t.global_order ASC").
		Limit(1).
		Offset(int64(order)).
		Bind(subQuery.Params()).
		WithContext(ctx).
		One(&cell)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return cellRowToCellInfo(&cell), nil
}

func (r *Repository) GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := baseCellQuery(pb.DB()).Build()

	var cells []cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		OrderBy("t.global_order ASC").
		Where(dbx.HashExp{
			pbhelper.DotExpand("t", schema.CellSchema.Type): string(t),
		}).
		Bind(subQuery.Params()).
		WithContext(ctx).
		All(&cells)
	if err != nil {
		return nil, err
	}

	return cellRowsToCellInfos(cells), nil
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

func (r *Repository) GetAllByWorldID(ctx context.Context, worldId string) ([]*model.CellInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	subQuery := baseCellQuery(pb.DB()).Build()

	var cells []cellRow
	err := pb.DB().
		Select("*").
		From(fmt.Sprintf("(%s) t", subQuery.SQL())).
		OrderBy("t.local_order ASC").
		Where(dbx.HashExp{
			pbhelper.DotExpand("t", schema.CellSchema.World): worldId,
		}).
		Bind(subQuery.Params()).
		WithContext(ctx).
		All(&cells)
	if err != nil {
		return nil, err
	}

	return cellRowsToCellInfos(cells), nil
}

func (r *Repository) GetAllIDs(ctx context.Context) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var ids []string
	err := pb.DB().
		Select(schema.CellSchema.Id).
		From(schema.CollectionCells).
		WithContext(ctx).
		Column(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) GetActivityFilterIDByID(ctx context.Context, id string) (string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var filterId string
	err := pb.DB().
		Select(schema.CellSchema.Filter).
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.Id: id,
		}).
		Limit(1).
		WithContext(ctx).
		Row(&filterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrCellNotFound
		}
		return "", err
	}

	return filterId, nil
}

func (r *Repository) UpdateAverageCampaignTimeByID(ctx context.Context, id string, campaignTime float64) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionCells,
			dbx.Params{
				schema.CellSchema.AverageCampaignTime: campaignTime,
			},
			dbx.HashExp{
				schema.CellSchema.Id: id,
			},
		).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrCellNotFound
	}

	return nil
}
