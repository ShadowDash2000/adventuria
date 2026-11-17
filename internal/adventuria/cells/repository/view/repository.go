package view

import (
	"adventuria/internal/adventuria/cells/service/view"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetAllBoard(ctx context.Context) ([]*view.CellBoard, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var cellRows []*cellBoardRow
	err := pb.DB().
		Select(
			schema.CellSchema.Id,
			schema.CellSchema.Type,
			schema.CellSchema.Icon,
			schema.CellSchema.World,
		).
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.Disabled: false,
		}).
		OrderBy(schema.CellSchema.Sort).
		WithContext(ctx).
		All(&cellRows)
	if err != nil {
		return nil, err
	}

	if len(cellRows) == 0 {
		return nil, nil
	}

	res := make([]*view.CellBoard, len(cellRows))
	for i, row := range cellRows {
		res[i] = row.toCellBoard()
	}

	return res, nil
}

func (r *Repository) GetDetailedByID(ctx context.Context, id string) (*view.CellDetailed, error) {
	cell, err := r.getCellDetailedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cell.Filter == "" {
		return cell.toCellDetailed(nil), nil
	}

	filterRow, err := r.getActivityFilterByID(ctx, cell.Filter)
	if err != nil {
		return nil, err
	}

	filter, err := r.getActivityFilterByFilterRow(ctx, filterRow)
	if err != nil {
		return nil, err
	}

	return cell.toCellDetailed(filter), nil
}

func (r *Repository) getCellDetailedByID(ctx context.Context, id string) (*cellDetailedRow, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var cell cellDetailedRow
	err := pb.DB().
		Select(
			schema.CellSchema.Id,
			schema.CellSchema.Type,
			schema.CellSchema.Name,
			schema.CellSchema.Icon,
			schema.CellSchema.Color,
			schema.CellSchema.Description,
			schema.CellSchema.Filter,
			schema.CellSchema.Points,
			schema.CellSchema.DoneEnergyConsume,
			schema.CellSchema.DropEnergyConsume,
			schema.CellSchema.Coins,
			schema.CellSchema.AverageCampaignTime,
		).
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.Id: id,
		}).
		WithContext(ctx).
		One(&cell)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCellNotFound
		}
		return nil, err
	}

	return &cell, err
}

func (r *Repository) getActivityFilterByID(ctx context.Context, id string) (*cellDetailedActivityFilterRow, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var filter *cellDetailedActivityFilterRow
	err := pb.DB().
		Select(
			schema.ActivityFilterSchema.Id,
			schema.ActivityFilterSchema.Type,
			schema.ActivityFilterSchema.Platforms,
			schema.ActivityFilterSchema.Developers,
			schema.ActivityFilterSchema.Publishers,
			schema.ActivityFilterSchema.Genres,
			schema.ActivityFilterSchema.Tags,
			schema.ActivityFilterSchema.Themes,
			schema.ActivityFilterSchema.MinPrice,
			schema.ActivityFilterSchema.MaxPrice,
			schema.ActivityFilterSchema.MinCampaignTime,
			schema.ActivityFilterSchema.MaxCampaignTime,
			schema.ActivityFilterSchema.Activities,
		).
		From(schema.CollectionActivityFilter).
		Where(dbx.HashExp{
			schema.ActivityFilterSchema.Id: id,
		}).
		WithContext(ctx).
		One(&filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityFilterNotFound
		}
		return nil, err
	}

	return filter, err
}

type activityFilterEntry[T any] interface {
	toDetailed() T
}

type activityFilterEntryColumns struct {
	id   string
	name string
}

func (r *Repository) getActivityFilterEntriesByIDs[E activityFilterEntry[T], T any](
	ctx context.Context,
	ids []string,
	table string,
	columns activityFilterEntryColumns,
) ([]T, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var entries []E
	err := pb.DB().
		Select(
			columns.id,
			columns.name,
		).
		From(table).
		Where(dbx.In(
			columns.id,
			pbhelper.SliceToAny(ids)...,
		)).
		WithContext(ctx).
		All(&entries)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, nil
	}

	res := make([]T, len(entries))
	for i, entry := range entries {
		res[i] = entry.toDetailed()
	}

	return res, nil
}

func (r *Repository) getActivitiesByIDs(ctx context.Context, ids []string) ([]*view.CellDetailedActivity, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var activities []*cellDetailedActivityRow
	err := pb.DB().
		Select(
			schema.ActivitySchema.Id,
			schema.ActivitySchema.Name,
			schema.ActivitySchema.HltbCampaignTime,
			schema.ActivitySchema.Cover,
			schema.ActivitySchema.CoverAlt,
		).
		From(schema.CollectionActivities).
		Where(dbx.In(
			schema.ActivitySchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		WithContext(ctx).
		All(&activities)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return nil, nil
	}

	res := make([]*view.CellDetailedActivity, len(activities))
	for i, activity := range activities {
		res[i] = activity.toCellDetailedActivity()
	}

	return res, nil
}

func (r *Repository) getActivityFilterByFilterRow(ctx context.Context, filterRow *cellDetailedActivityFilterRow) (*view.CellDetailedActivityFilter, error) {
	platforms, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedPlatformRow,
		*view.CellDetailedPlatform,
	](ctx, filterRow.Platforms, schema.CollectionPlatforms, activityFilterEntryColumns{
		id:   schema.PlatformSchema.Id,
		name: schema.PlatformSchema.Name,
	})
	if err != nil {
		return nil, err
	}

	developers, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedCompanyRow,
		*view.CellDetailedCompany,
	](ctx, filterRow.Developers, schema.CollectionCompanies, activityFilterEntryColumns{
		id:   schema.CompanySchema.Id,
		name: schema.CompanySchema.Name,
	})
	if err != nil {
		return nil, err
	}

	publishers, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedCompanyRow,
		*view.CellDetailedCompany,
	](ctx, filterRow.Publishers, schema.CollectionCompanies, activityFilterEntryColumns{
		id:   schema.CompanySchema.Id,
		name: schema.CompanySchema.Name,
	})
	if err != nil {
		return nil, err
	}

	genres, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedGenreRow,
		*view.CellDetailedGenre,
	](ctx, filterRow.Genres, schema.CollectionGenres, activityFilterEntryColumns{
		id:   schema.GenreSchema.Id,
		name: schema.GenreSchema.Name,
	})
	if err != nil {
		return nil, err
	}

	tags, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedTagRow,
		*view.CellDetailedTag,
	](ctx, filterRow.Tags, schema.CollectionTags, activityFilterEntryColumns{
		id:   schema.TagSchema.Id,
		name: schema.TagSchema.Name,
	})
	if err != nil {
		return nil, err
	}

	themes, err := r.getActivityFilterEntriesByIDs[
		*cellDetailedThemeRow,
		*view.CellDetailedTheme,
	](ctx, filterRow.Themes, schema.CollectionThemes, activityFilterEntryColumns{
		id:   schema.ThemeSchema.Id,
		name: schema.ThemeSchema.Name,
	})
	if err != nil {
		return nil, err
	}

	activities, err := r.getActivitiesByIDs(ctx, filterRow.Activities)
	if err != nil {
		return nil, err
	}

	return &view.CellDetailedActivityFilter{
		Type:            filterRow.Type,
		Platforms:       platforms,
		Developers:      developers,
		Publishers:      publishers,
		Genres:          genres,
		Tags:            tags,
		Themes:          themes,
		MinPrice:        filterRow.MinPrice,
		MaxPrice:        filterRow.MaxPrice,
		MinCampaignTime: filterRow.MinCampaignTime,
		MaxCampaignTime: filterRow.MaxCampaignTime,
		Activities:      activities,
	}, nil
}
