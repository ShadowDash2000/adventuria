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

func (r *Repository) Create(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionActivities)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	ActivityToRecord(activity, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToActivity(record), nil
}

func (r *Repository) Update(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ActivitySchema.Id: activity.ID(),
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityNotFound
		}
		return nil, err
	}

	ActivityToRecord(activity, &record)
	err = pb.SaveWithContext(ctx, &record)
	if err != nil {
		return nil, err
	}

	return RecordToActivity(&record), nil
}

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ActivitySchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityNotFound
		}

		return nil, err
	}

	return RecordToActivity(&record), nil
}

func (r *Repository) GetCountByFilter(ctx context.Context, filter model.ActivityFilter) (int, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	countQuery := pb.DB().Select("count(*)")
	buildQuery(pb, filter, countQuery)

	var totalCount int
	err := countQuery.WithContext(ctx).Row(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *Repository) GetIDsByFilter(ctx context.Context, filter model.ActivityFilter, offset, limit int) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	q := pb.DB().
		Select(fmt.Sprintf("f.%s", schema.ActivitySchema.Id)).
		Limit(int64(limit)).
		Offset(int64(offset))
	buildQuery(pb, filter, q, schema.ActivitySchema.Id)

	var ids []string
	err := q.WithContext(ctx).Column(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) GetAverageCampaignTimeByFilter(ctx context.Context, filter model.ActivityFilter) (float64, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	q := pb.DB().
		Select(
			fmt.Sprintf(
				"avg(nullif(f.%s, 0))",
				schema.ActivitySchema.HltbCampaignTime,
			),
		)
	buildQuery(
		pb, filter, q,
		schema.ActivitySchema.Id,
		schema.ActivitySchema.HltbCampaignTime,
	)

	var campaignTime sql.NullFloat64
	err := q.WithContext(ctx).Row(&campaignTime)
	if err != nil {
		return 0, err
	}

	if !campaignTime.Valid {
		return 0, nil
	}

	return campaignTime.Float64, nil
}

func buildQuery(app core.App, filter model.ActivityFilter, mainQuery *dbx.SelectQuery, fields ...string) {
	q := app.DB().
		Select(fields...).
		From(schema.CollectionActivities)

	// if ids are specified, then we don't need any other filters
	if len(filter.Activities) > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf(
				"%s IN (%s)",
				schema.ActivitySchema.Id,
				pbhelper.SliceToSqlString(filter.Activities),
			),
		))

		mainQuery.From(fmt.Sprintf("(%s) AS f", q.Build().SQL()))

		return
	}

	if filter.Type != "" {
		q.Where(dbx.NewExp(fmt.Sprintf("%s = '%s'", schema.ActivitySchema.Type, filter.Type)))
	}

	if len(filter.GameTypes) > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf(
				"%s IN (%s)",
				schema.ActivitySchema.GameType,
				pbhelper.SliceToSqlString(filter.GameTypes),
			),
		))
	}

	if filter.MinPrice > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > %d", schema.ActivitySchema.SteamAppPrice, filter.MinPrice),
		))
	}
	if filter.MaxPrice > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < %d", schema.ActivitySchema.SteamAppPrice, filter.MaxPrice),
		))
	}

	if !filter.ReleaseDateFrom.IsZero() {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > '%s'", schema.ActivitySchema.ReleaseDate, filter.ReleaseDateFrom.String()),
		))
	}
	if !filter.ReleaseDateTo.IsZero() {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < '%s'", schema.ActivitySchema.ReleaseDate, filter.ReleaseDateTo.String()),
		))
	}

	if filter.MinCampaignTime > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > %f", schema.ActivitySchema.HltbCampaignTime, filter.MinCampaignTime),
		))
	}
	if filter.MaxCampaignTime > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < %f", schema.ActivitySchema.HltbCampaignTime, filter.MaxCampaignTime),
		))
	}

	mainQuery.From(fmt.Sprintf("(%s) AS f", q.Build().SQL()))

	setSubTablesFilters(app, filter, mainQuery)
}

func setSubTablesFilters(app core.App, filter model.ActivityFilter, q *dbx.SelectQuery) {
	if len(filter.Platforms) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesPlatforms,
			schema.ActivitiesPlatformsSchema.Activity,
			schema.ActivitiesPlatformsSchema.Platform,
			filter.Platforms,
			filter.PlatformsStrict,
		)
	}
	if len(filter.Developers) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesDevelopers,
			schema.ActivitiesDevelopersSchema.Activity,
			schema.ActivitiesDevelopersSchema.Developer,
			filter.Developers,
			false,
		)
	}
	if len(filter.Publishers) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesPublishers,
			schema.ActivitiesPublishersSchema.Activity,
			schema.ActivitiesPublishersSchema.Publisher,
			filter.Publishers,
			false,
		)
	}
	if len(filter.Genres) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesGenres,
			schema.ActivitiesGenresSchema.Activity,
			schema.ActivitiesGenresSchema.Genre,
			filter.Genres,
			false,
		)
	}
	if len(filter.Tags) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesTags,
			schema.ActivitiesTagsSchema.Activity,
			schema.ActivitiesTagsSchema.Tag,
			filter.Tags,
			false,
		)
	}
	if len(filter.Themes) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesThemes,
			schema.ActivitiesThemesSchema.Activity,
			schema.ActivitiesThemesSchema.Theme,
			filter.Themes,
			false,
		)
	}
}

func applyActivityRelationFilter(
	pb core.App,
	query *dbx.SelectQuery,
	collectionName string,
	activityField string,
	relationField string,
	values []string,
	strict bool,
) {
	if len(values) == 0 {
		return
	}

	inClause := pbhelper.SliceToSqlString(values)

	subQuery := pb.DB().
		Select(activityField).
		From(collectionName).
		Where(dbx.NewExp(fmt.Sprintf("%s IN (%s)", relationField, inClause))).
		Build()

	query.AndWhere(dbx.NewExp(fmt.Sprintf("id IN (%s)", subQuery.SQL())))

	if strict {
		mainIdField := fmt.Sprintf("%s.id", "f")
		subQuery := pb.DB().
			Select(activityField).
			From(collectionName).
			GroupBy(activityField).
			Having(dbx.NewExp("COUNT(*) = 1")).
			Build()

		query.AndWhere(
			dbx.NewExp(
				fmt.Sprintf("%s IN (%s)", mainIdField, subQuery.SQL()),
			),
		)
	}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.HashExp{schema.ActivitySchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityNotFound
		}
		return nil, err
	}

	return RecordToActivity(&record), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error) {
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

	return RecordsToActivities(records), nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ActivitySchema.Name: name,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityNotFound
		}
		return nil, err
	}

	return RecordToActivity(&record), nil
}

func (r *Repository) GetByNameParts(ctx context.Context, nameParts []string) ([]*model.Activity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Where(dbx.Like(schema.ActivitySchema.Name, nameParts...)).
		All(&records)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errs.ErrActivityNotFound
	}

	return RecordsToActivities(records), nil
}
