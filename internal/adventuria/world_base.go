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
	return w.GetString(schema.WorldsSchema.Name)
}

func (w *WorldBase) Slug() string {
	return w.GetString(schema.WorldsSchema.Slug)
}

func (w *WorldBase) Sort() int {
	return w.GetInt(schema.WorldsSchema.Sort)
}

func (w *WorldBase) IsLoop() bool {
	return w.GetBool(schema.WorldsSchema.IsLoop)
}

func (w *WorldBase) IsDefaultWorld() bool {
	return w.GetBool(schema.WorldsSchema.IsDefaultWorld)
}

func (w *WorldBase) TransitionToWorld() string {
	return w.GetString(schema.WorldsSchema.TransitionToWorld)
}

func (w *WorldBase) Effects() []string {
	return w.GetStringSlice(schema.WorldsSchema.Effects)
}
