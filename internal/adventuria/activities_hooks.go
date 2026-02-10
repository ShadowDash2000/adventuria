package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func BindActivitiesHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionActivities).BindFunc(OnAfterActivityCreateSuccess)
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionActivities).BindFunc(OnAfterActivityUpdateSuccess)
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionActivities).BindFunc(OnAfterActivityDeleteSuccess)
}

func OnAfterActivityCreateSuccess(e *core.RecordEvent) error {
	err := e.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}
		for field, collection := range indexCollections {
			err := collection.createRecords(ctx, e.Record.Id, e.Record.GetStringSlice(field))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return e.Next()
}

func OnAfterActivityUpdateSuccess(e *core.RecordEvent) error {
	err := e.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}
		for field, collection := range indexCollections {
			err := collection.updateRecords(ctx, e.Record.Id, e.Record.GetStringSlice(field))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return e.Next()
}

func OnAfterActivityDeleteSuccess(e *core.RecordEvent) error {
	err := e.App.RunInTransaction(func(txApp core.App) error {
		ctx := AppContext{App: txApp}
		for _, collection := range indexCollections {
			err := collection.deleteRecords(ctx, e.Record.Id)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return e.Next()
}

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

func (i *indexCollection) createRecords(ctx AppContext, activityId string, relationIds []string) error {
	for _, relationId := range relationIds {
		err := i.createRecord(ctx, activityId, relationId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *indexCollection) createRecord(ctx AppContext, activityId, relationId string) error {
	record := core.NewRecord(GameCollections.Get(i.collection))
	record.Set(i.activityField, activityId)
	record.Set(i.relationField, relationId)
	return ctx.App.Save(record)
}

func (i *indexCollection) updateRecords(ctx AppContext, activityId string, relationIds []string) error {
	var records []*core.Record
	err := ctx.App.
		RecordQuery(i.collection).
		Where(dbx.HashExp{i.activityField: activityId}).
		All(&records)
	if err != nil {
		return err
	}

	toDelete := make(map[string]*core.Record, len(records))
	for _, record := range records {
		toDelete[record.Id] = record
	}

	for _, relationId := range relationIds {
		if _, ok := toDelete[relationId]; ok {
			delete(toDelete, relationId)
		} else {
			err = i.createRecord(ctx, activityId, relationId)
			if err != nil {
				return err
			}
		}
	}

	for _, record := range toDelete {
		err = ctx.App.Delete(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *indexCollection) deleteRecords(ctx AppContext, activityId string) error {
	var records []*core.Record
	err := ctx.App.
		RecordQuery(i.collection).
		Where(dbx.HashExp{i.activityField: activityId}).
		All(&records)
	if err != nil {
		return err
	}

	for _, record := range records {
		err = ctx.App.Delete(record)
		if err != nil {
			return err
		}
	}

	return nil
}
