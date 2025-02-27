package adventuria

import "github.com/pocketbase/pocketbase/core"

type Log struct {
	app           core.App
	logCollection *core.Collection
}

func NewLog(app core.App) *Log {
	return &Log{app: app}
}

func (l *Log) Init() error {
	var err error
	l.logCollection, err = l.app.FindCollectionByNameOrId(TableLogs)
	if err != nil {
		return err
	}
	return nil
}

func (l *Log) Add(userId, logType, value string) error {
	record := core.NewRecord(l.logCollection)
	record.Set("user", userId)
	record.Set("type", logType)
	record.Set("value", value)

	err := l.app.Save(record)
	if err != nil {
		return err
	}

	return nil
}
