package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

type WorldBase struct {
	core.BaseRecordProxy
}

func NewWorld(record *core.Record) World {
	w := &WorldBase{}
	w.SetProxyRecord(record)
	return w
}

func (w *WorldBase) ID() string {
	return w.Id
}

func (w *WorldBase) Name() string {
	return w.GetString(schema.WorldSchema.Name)
}

func (w *WorldBase) Slug() string {
	return w.GetString(schema.WorldSchema.Slug)
}

func (w *WorldBase) Sort() int {
	return w.GetInt(schema.WorldSchema.Sort)
}

func (w *WorldBase) IsLoop() bool {
	return w.GetBool(schema.WorldSchema.IsLoop)
}

func (w *WorldBase) IsDefaultWorld() bool {
	return w.GetBool(schema.WorldSchema.IsDefaultWorld)
}

func (w *WorldBase) TransitionToWorld() string {
	return w.GetString(schema.WorldSchema.TransitionToWorld)
}

func (w *WorldBase) Effects() []string {
	return w.GetStringSlice(schema.WorldSchema.Effects)
}
