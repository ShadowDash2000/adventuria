package activities

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetOrCreate(ctx context.Context, data model.ActivityCreate) (*model.Activity, error)
	GetActivitiesByFilter(ctx context.Context, filter *model.ActivityFilter, poolSize, resultSize int) ([]string, error)
	GetByID(ctx context.Context, id string) (*model.Activity, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, activity *model.Activity) (*model.Activity, error)
}

type Activities struct {
	repository repository
}

func NewActivities(repository repository) *Activities {
	return &Activities{repository: repository}
}

func (a *Activities) GetOrCreate(ctx context.Context, data model.ActivityCreate) (*model.Activity, error) {
	return a.repository.GetOrCreate(ctx, data)
}

func (a *Activities) UpdateActivitiesFromFilter(
	ctx context.Context,
	player *model.Player,
	filter *model.ActivityFilter,
	forceUpdate bool,
) error {
	needToUpdate := forceUpdate
	customFilter := player.LastAction().CustomActivityFilter()

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
		activitiesIds, err := a.GetByFilter(ctx, filter)
		if err != nil {
			return err
		}

		player.LastAction().SetItemsList(activitiesIds)
	}

	return nil
}

func (a *Activities) GetByFilter(ctx context.Context, filter *model.ActivityFilter) ([]string, error) {
	return a.repository.GetActivitiesByFilter(ctx, filter, 20000, 20)
}

func (a *Activities) GetByID(ctx context.Context, id string) (*model.Activity, error) {
	return a.repository.GetByID(ctx, id)
}

func (a *Activities) GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error) {
	return a.repository.GetByIDs(ctx, ids)
}

func (a *Activities) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return a.repository.GetChecksumsByIDs(ctx, ids)
}

func (a *Activities) Save(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	return a.repository.Save(ctx, activity)
}
