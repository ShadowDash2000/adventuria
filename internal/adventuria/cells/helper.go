package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func newActivityFilterById(app core.App, filterId string) (adventuria.ActivityFilterRecord, error) {
	var filter adventuria.ActivityFilterRecord
	if filterId != "" {
		filterRecord, err := app.FindRecordById(
			adventuria.GameCollections.Get(schema.CollectionActivityFilter),
			filterId,
		)
		if err != nil {
			return nil, err
		}

		filter = adventuria.NewActivityFilterFromRecord(filterRecord)
	} else {
		filter = adventuria.NewActivityFilterFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(schema.CollectionActivityFilter),
			),
		)
	}
	return filter, nil
}

func updateActivitiesFromFilter(app core.App, user adventuria.User, filter adventuria.ActivityFilterRecord, forceUpdate bool) error {
	needToUpdate := forceUpdate
	customFilter, err := user.LastAction().CustomActivityFilter()
	if err != nil {
		return err
	}

	if len(customFilter.Platforms) > 0 {
		filter.SetPlatforms(append(filter.Platforms(), customFilter.Platforms...))
		needToUpdate = true
	}
	if len(customFilter.Developers) > 0 {
		filter.SetDevelopers(append(filter.Developers(), customFilter.Developers...))
		needToUpdate = true
	}
	if len(customFilter.Publishers) > 0 {
		filter.SetPublishers(append(filter.Publishers(), customFilter.Publishers...))
		needToUpdate = true
	}
	if len(customFilter.Genres) > 0 {
		filter.SetGenres(append(filter.Genres(), customFilter.Genres...))
		needToUpdate = true
	}
	if len(customFilter.Tags) > 0 {
		filter.SetTags(append(filter.Tags(), customFilter.Tags...))
		needToUpdate = true
	}
	if len(customFilter.Themes) > 0 {
		filter.SetThemes(append(filter.Themes(), customFilter.Themes...))
		needToUpdate = true
	}
	if customFilter.MinPrice != 0 {
		filter.SetMinPrice(customFilter.MinPrice)
		needToUpdate = true
	}
	if customFilter.MaxPrice != 0 {
		filter.SetMaxPrice(customFilter.MaxPrice)
		needToUpdate = true
	}
	if !customFilter.ReleaseDateFrom.IsZero() {
		filter.SetReleaseDateFrom(customFilter.ReleaseDateFrom)
		needToUpdate = true
	}
	if !customFilter.ReleaseDateTo.IsZero() {
		filter.SetReleaseDateTo(customFilter.ReleaseDateTo)
		needToUpdate = true
	}
	if customFilter.MinCampaignTime != 0 {
		filter.SetMinCampaignTime(customFilter.MinCampaignTime)
		needToUpdate = true
	}
	if customFilter.MaxCampaignTime != 0 {
		filter.SetMaxCampaignTime(customFilter.MaxCampaignTime)
		needToUpdate = true
	}

	if needToUpdate {
		res, err := fetchActivitiesByFilter(app, filter)
		if err != nil {
			return err
		}

		user.LastAction().SetItemsList(res)
	}

	return nil
}

func fetchActivitiesByFilter(app core.App, filter adventuria.ActivityFilterRecord) ([]string, error) {
	countQuery := app.RecordQuery(adventuria.GameCollections.Get(schema.CollectionActivities))
	if filter != nil {
		setFilters(app, filter, countQuery)
	}

	var totalCount int
	err := countQuery.Select("count(activities.id)").Row(&totalCount)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return []string{}, nil
	}

	const maxPoolSize = 20000
	limit := totalCount
	offset := 0

	if totalCount > maxPoolSize {
		limit = maxPoolSize
		offset = rand.N(totalCount - maxPoolSize + 1)
	}

	q := app.RecordQuery(adventuria.GameCollections.Get(schema.CollectionActivities)).
		Select("activities.id").
		Limit(int64(limit)).
		Offset(int64(offset))
	if filter != nil {
		setFilters(app, filter, q)
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = q.All(&records)
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})

	resultLimit := 20
	if len(records) < resultLimit {
		resultLimit = len(records)
	}

	res := make([]string, resultLimit)
	for i := 0; i < resultLimit; i++ {
		res[i] = records[i].Id
	}

	return res, nil
}

func setFilters(app core.App, filter adventuria.ActivityFilterRecord, q *dbx.SelectQuery) {
	if filter.Type() != "" {
		q = q.Where(dbx.NewExp(fmt.Sprintf("%s = '%s'", schema.ActivitySchema.Type, filter.Type())))
	}

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
	if len(filter.GameTypes()) > 0 {
		q.AndWhere(dbx.In(schema.ActivitySchema.GameType, sliceToAny(filter.GameTypes())...))
	}
	if len(filter.Activities()) > 0 {
		q.AndWhere(dbx.In("id", sliceToAny(filter.Activities())...))
	}

	if filter.MinPrice() > 0 {
		q.AndWhere(dbx.NewExp("steam_app_price > {:price}", dbx.Params{"price": filter.MinPrice()}))
	}
	if filter.MaxPrice() > 0 {
		q.AndWhere(dbx.NewExp("steam_app_price < {:price}", dbx.Params{"price": filter.MaxPrice()}))
	}

	if !filter.ReleaseDateFrom().IsZero() {
		q.AndWhere(dbx.NewExp("release_date > {:date}", dbx.Params{"date": filter.ReleaseDateFrom()}))
	}
	if !filter.ReleaseDateTo().IsZero() {
		q.AndWhere(dbx.NewExp("release_date < {:date}", dbx.Params{"date": filter.ReleaseDateTo()}))
	}

	if filter.MinCampaignTime() > 0 {
		q.AndWhere(dbx.NewExp("hltb_campaign_time > {:min_time}", dbx.Params{"min_time": filter.MinCampaignTime()}))
	}
	if filter.MaxCampaignTime() > 0 {
		q.AndWhere(dbx.NewExp("hltb_campaign_time < {:max_time}", dbx.Params{"max_time": filter.MaxCampaignTime()}))
	}
}

func sliceToAny[T any](slice []T) []any {
	res := make([]any, len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res
}

func applyActivityRelationFilter(
	app core.App,
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

	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("'%s'", v)
	}
	inClause := strings.Join(quotedValues, ", ")

	subQuery := app.DB().
		Select(activityField).
		From(collectionName).
		Where(dbx.NewExp(fmt.Sprintf("%s IN (%s)", relationField, inClause))).
		Build()

	query.AndWhere(dbx.NewExp(fmt.Sprintf("id IN (%s)", subQuery.SQL())))

	if strict {
		mainIdField := fmt.Sprintf("%s.id", schema.CollectionActivities)
		subQuery := app.DB().
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
