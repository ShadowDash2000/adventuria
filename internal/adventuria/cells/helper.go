package cells

import (
	"adventuria/internal/adventuria"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func newActivityFilterById(filterId string) (adventuria.ActivityFilterRecord, error) {
	var filter adventuria.ActivityFilterRecord
	if filterId != "" {
		filterRecord, err := adventuria.PocketBase.FindRecordById(
			adventuria.GameCollections.Get(adventuria.CollectionActivityFilter),
			filterId,
		)
		if err != nil {
			return nil, err
		}

		filter = adventuria.NewActivityFilterFromRecord(filterRecord)
	} else {
		filter = adventuria.NewActivityFilterFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(adventuria.CollectionActivityFilter),
			),
		)
	}
	return filter, nil
}

func updateActivitiesFromFilter(user adventuria.User, filter adventuria.ActivityFilterRecord, forceUpdate bool) error {
	needToUpdate := forceUpdate
	customFilter := user.LastAction().CustomActivityFilter()

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
		res, err := fetchActivitiesByFilter(filter)
		if err != nil {
			return err
		}

		user.LastAction().SetItemsList(res)
	}

	return nil
}

func fetchActivitiesByFilter(filter adventuria.ActivityFilterRecord) ([]string, error) {
	q := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionActivities)).
		Limit(20).
		OrderBy("random()")

	if filter != nil {
		q = setFilters(filter, q)
	}

	var records []*core.Record
	err := q.All(&records)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(records))
	for i, record := range records {
		res[i] = record.Id
	}

	return res, nil
}

func setFilters(filter adventuria.ActivityFilterRecord, q *dbx.SelectQuery) *dbx.SelectQuery {
	if filter.Type() != "" {
		q = q.AndWhere(dbx.HashExp{"type": filter.Type()})
	}

	if len(filter.Platforms()) > 0 {
		if filter.PlatformsStrict() {
			exps := make([]dbx.Expression, len(filter.Platforms()))
			for i, platform := range filter.Platforms() {
				exps[i] = dbx.Like("platforms", fmt.Sprintf("[\"%s\"]", platform))
			}
			q = q.AndWhere(dbx.Or(exps...))
		} else {
			q = q.AndWhere(dbx.OrLike("platforms", filter.Platforms()...))
		}
	}
	if len(filter.Developers()) > 0 {
		q = q.AndWhere(dbx.OrLike("developers", filter.Developers()...))
	}
	if len(filter.Publishers()) > 0 {
		q = q.AndWhere(dbx.OrLike("publishers", filter.Publishers()...))
	}
	if len(filter.Genres()) > 0 {
		q = q.AndWhere(dbx.OrLike("genres", filter.Genres()...))
	}
	if len(filter.Tags()) > 0 {
		q = q.AndWhere(dbx.OrLike("tags", filter.Tags()...))
	}
	if len(filter.Themes()) > 0 {
		q = q.AndWhere(dbx.OrLike("themes", filter.Themes()...))
	}
	if len(filter.Games()) > 0 {
		q = q.AndWhere(dbx.OrLike("id", filter.Games()...))
	}

	if filter.MinPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_app_price > {:price}", dbx.Params{"price": filter.MinPrice()}))
	}
	if filter.MaxPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_app_price < {:price}", dbx.Params{"price": filter.MaxPrice()}))
	}

	if !filter.ReleaseDateFrom().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date > {:date}", dbx.Params{"date": filter.ReleaseDateFrom()}))
	}
	if !filter.ReleaseDateTo().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date < {:date}", dbx.Params{"date": filter.ReleaseDateTo()}))
	}

	if filter.MinCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("hltb_campaign_time > {:time}", dbx.Params{"time": filter.MinCampaignTime()}))
	}
	if filter.MaxCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("hltb_campaign_time < {:time}", dbx.Params{"time": filter.MaxCampaignTime()}))
	}

	return q
}
