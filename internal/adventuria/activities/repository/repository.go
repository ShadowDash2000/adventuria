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
	"math/rand/v2"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
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

func (r *Repository) GetActivitiesByFilter(ctx context.Context, filter *model.ActivityFilter, poolSize, resultSize int) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	countQuery := pb.DB().Select("count(*)")
	if filter != nil {
		buildQuery(pb, filter, countQuery)
	}

	var totalCount int
	err := countQuery.WithContext(ctx).Row(&totalCount)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return []string{}, nil
	}

	limit := totalCount
	offset := 0

	if totalCount > poolSize {
		limit = poolSize
		offset = rand.N(totalCount - poolSize + 1)
	}

	q := pb.DB().
		Select("f.id").
		Limit(int64(limit)).
		Offset(int64(offset))
	if filter != nil {
		buildQuery(pb, filter, q)
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = q.WithContext(ctx).All(&records)
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})

	if len(records) < resultSize {
		resultSize = len(records)
	}

	res := make([]string, resultSize)
	for i := 0; i < resultSize; i++ {
		res[i] = records[i].Id
	}

	return res, nil
}

func buildQuery(app core.App, filter *model.ActivityFilter, mainQuery *dbx.SelectQuery) {
	q := app.DB().
		Select("id").
		From(schema.CollectionActivities)

	// if ids are specified, then we don't need any other filters
	if len(filter.Activities()) > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf(
				"%s IN (%s)",
				schema.ActivitySchema.Id,
				pbhelper.SliceToSqlString(filter.Activities()),
			),
		))

		mainQuery.From(fmt.Sprintf("(%s) AS f", q.Build().SQL()))

		return
	}

	if filter.Type() != "" {
		q.Where(dbx.NewExp(fmt.Sprintf("%s = '%s'", schema.ActivitySchema.Type, filter.Type())))
	}

	if len(filter.GameTypes()) > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf(
				"%s IN (%s)",
				schema.ActivitySchema.GameType,
				pbhelper.SliceToSqlString(filter.GameTypes()),
			),
		))
	}

	if filter.MinPrice() > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > %d", schema.ActivitySchema.SteamAppPrice, filter.MinPrice()),
		))
	}
	if filter.MaxPrice() > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < %d", schema.ActivitySchema.SteamAppPrice, filter.MaxPrice()),
		))
	}

	if !filter.ReleaseDateFrom().IsZero() {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > '%s'", schema.ActivitySchema.ReleaseDate, filter.ReleaseDateFrom().String()),
		))
	}
	if !filter.ReleaseDateTo().IsZero() {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < '%s'", schema.ActivitySchema.ReleaseDate, filter.ReleaseDateTo().String()),
		))
	}

	if filter.MinCampaignTime() > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > %f", schema.ActivitySchema.HltbCampaignTime, filter.MinCampaignTime()),
		))
	}
	if filter.MaxCampaignTime() > 0 {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s < %f", schema.ActivitySchema.HltbCampaignTime, filter.MaxCampaignTime()),
		))
	}

	mainQuery.From(fmt.Sprintf("(%s) AS f", q.Build().SQL()))

	setSubTablesFilters(app, filter, mainQuery)
}

func setSubTablesFilters(app core.App, filter *model.ActivityFilter, q *dbx.SelectQuery) {
	if len(filter.Platforms()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesPlatforms,
			schema.ActivitiesPlatformsSchema.Activity,
			schema.ActivitiesPlatformsSchema.Platform,
			filter.Platforms(),
			filter.PlatformsStrict(),
		)
	}
	if len(filter.Developers()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesDevelopers,
			schema.ActivitiesDevelopersSchema.Activity,
			schema.ActivitiesDevelopersSchema.Developer,
			filter.Developers(),
			false,
		)
	}
	if len(filter.Publishers()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesPublishers,
			schema.ActivitiesPublishersSchema.Activity,
			schema.ActivitiesPublishersSchema.Publisher,
			filter.Publishers(),
			false,
		)
	}
	if len(filter.Genres()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesGenres,
			schema.ActivitiesGenresSchema.Activity,
			schema.ActivitiesGenresSchema.Genre,
			filter.Genres(),
			false,
		)
	}
	if len(filter.Tags()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesTags,
			schema.ActivitiesTagsSchema.Activity,
			schema.ActivitiesTagsSchema.Tag,
			filter.Tags(),
			false,
		)
	}
	if len(filter.Themes()) > 0 {
		applyActivityRelationFilter(
			app,
			q,
			schema.CollectionActivitiesThemes,
			schema.ActivitiesThemesSchema.Activity,
			schema.ActivitiesThemesSchema.Theme,
			filter.Themes(),
			false,
		)
	}
}

func applyActivityRelationFilter(
	pb core.App,
	query *dbx.SelectQuery,
	collectionName,
	activityField,
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

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionActivities).
		WithContext(ctx).
		Select(
			schema.ActivitySchema.Id,
			schema.ActivitySchema.Checksum,
		).
		Where(dbx.In(
			schema.ActivitySchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.ActivitySchema.Checksum)
	}

	return checksums, nil
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

	record, err := pb.FindRecordById(schema.CollectionActivities, activity.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActivityNotFound
		}
		return nil, err
	}

	ActivityToRecord(activity, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToActivity(record), nil
}

func (r *Repository) Save(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	if activity.IsNew() {
		return r.Create(ctx, activity)
	}

	return r.Update(ctx, activity)
}
