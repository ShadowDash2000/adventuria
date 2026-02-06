package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

type Timers struct{}

func NewTimers(ctx AppContext) *Timers {
	t := &Timers{}
	t.bindHooks(ctx)
	return t
}

func (t *Timers) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionTimers).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString(schema.TimerSchema.User)

		user, err := GameUsers.GetByID(AppContext{App: e.App}, userId)
		if err != nil {
			return e.Next()
		}

		user.Timer().SetProxyRecord(e.Record)

		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionTimers).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString(schema.TimerSchema.User)

		user, err := GameUsers.GetByID(AppContext{App: e.App}, userId)
		if err != nil {
			return e.Next()
		}

		user.Timer().SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionTimers)))

		return e.Next()
	})
}
