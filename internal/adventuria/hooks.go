package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func (g *Game) bindHooks(pb core.App) {
	pb.OnRecordUpdate(schema.CollectionSettings).BindFunc(func(e *core.RecordEvent) error {
		if ok := e.Record.GetBool(schema.SettingsSchema.KillParser); ok {
			err := g.onKillParser().Trigger(e.Context, &onKillParserEvent{})
			if err != nil {
				e.App.Logger().Error("Failed to trigger kill parser event", "err", err)
			}
			e.Record.Set(schema.SettingsSchema.KillParser, false)
		}
		return e.Next()
	})
}
