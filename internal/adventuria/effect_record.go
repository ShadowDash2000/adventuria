package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type EffectRecord struct {
	core.BaseRecordProxy
}

func (e *EffectRecord) ID() string {
	return e.Id
}

func (e *EffectRecord) Name() string {
	return e.GetString("name")
}

func (e *EffectRecord) Type() string {
	return e.GetString("type")
}
