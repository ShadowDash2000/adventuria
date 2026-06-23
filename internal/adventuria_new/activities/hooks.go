package activities

import (
	"adventuria/internal/adventuria/schema"
	repo "adventuria/internal/adventuria_new/activities/repository"
	"adventuria/pkg/helper"
	"adventuria/pkg/pbtransaction"
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type indexCollection struct {
	collection    string
	activityField string
	relationField string
}

var indexCollections = map[string]indexCollection{
	schema.ActivitySchema.Platforms: {
		collection:    schema.CollectionActivitiesPlatforms,
		activityField: schema.ActivitiesPlatformsSchema.Activity,
		relationField: schema.ActivitiesPlatformsSchema.Platform,
	},
	schema.ActivitySchema.Developers: {
		collection:    schema.CollectionActivitiesDevelopers,
		activityField: schema.ActivitiesDevelopersSchema.Activity,
		relationField: schema.ActivitiesDevelopersSchema.Developer,
	},
	schema.ActivitySchema.Publishers: {
		collection:    schema.CollectionActivitiesPublishers,
		activityField: schema.ActivitiesPublishersSchema.Activity,
		relationField: schema.ActivitiesPublishersSchema.Publisher,
	},
	schema.ActivitySchema.Genres: {
		collection:    schema.CollectionActivitiesGenres,
		activityField: schema.ActivitiesGenresSchema.Activity,
		relationField: schema.ActivitiesGenresSchema.Genre,
	},
	schema.ActivitySchema.Tags: {
		collection:    schema.CollectionActivitiesTags,
		activityField: schema.ActivitiesTagsSchema.Activity,
		relationField: schema.ActivitiesTagsSchema.Tag,
	},
	schema.ActivitySchema.Themes: {
		collection:    schema.CollectionActivitiesThemes,
		activityField: schema.ActivitiesThemesSchema.Activity,
		relationField: schema.ActivitiesThemesSchema.Theme,
	},
}

func BindHooks(app core.App, repo *repo.RelationRepository) {
	app.OnRecordAfterCreateSuccess(schema.CollectionActivities).BindFunc(func(e *core.RecordEvent) error {
		return pbtransaction.RunInTransaction(e.Context, e.App, func(ctx context.Context, txApp core.App) error {
			for field, ic := range indexCollections {
				relationIds := e.Record.GetStringSlice(field)
				err := repo.SyncRelations(e.Context, ic.collection, ic.activityField, ic.relationField, e.Record.Id, relationIds)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})

	app.OnRecordAfterUpdateSuccess(schema.CollectionActivities).BindFunc(func(e *core.RecordEvent) error {
		original := e.Record.Original()
		return pbtransaction.RunInTransaction(e.Context, e.App, func(ctx context.Context, txApp core.App) error {
			for field, ic := range indexCollections {
				currentSlice := e.Record.GetStringSlice(field)
				originalSlice := original.GetStringSlice(field)

				changed := len(currentSlice) != len(originalSlice) || !helper.SliceContainsAll(currentSlice, originalSlice)

				if !changed {
					continue
				}

				err := repo.SyncRelations(e.Context, ic.collection, ic.activityField, ic.relationField, e.Record.Id, currentSlice)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})

	app.OnRecordAfterDeleteSuccess(schema.CollectionActivities).BindFunc(func(e *core.RecordEvent) error {
		return pbtransaction.RunInTransaction(e.Context, e.App, func(ctx context.Context, txApp core.App) error {
			for _, ic := range indexCollections {
				err := repo.DeleteAllRelations(e.Context, ic.collection, ic.activityField, e.Record.Id)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
}
