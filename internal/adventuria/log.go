package adventuria

import (
	"adventuria/pkg/collections"
	"github.com/pocketbase/pocketbase/core"
)

const (
	LogTypeItemUse           = "itemUse"
	LogTypeItemDrop          = "itemDrop"
	LogTypeItemEffectApplied = "itemEffectApplied"
)

type Log struct {
	app  core.App
	cols *collections.Collections
}

func NewLog(cols *collections.Collections, app core.App) *Log {
	return &Log{
		app:  app,
		cols: cols,
	}
}

func (l *Log) Add(userId, logType, value string) error {
	collection, err := l.cols.Get(TableLogs)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Set("user", userId)
	record.Set("type", logType)
	record.Set("value", value)

	err = l.app.Save(record)
	if err != nil {
		return err
	}

	return nil
}
