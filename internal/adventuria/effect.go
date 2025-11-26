package adventuria

import (
	"adventuria/pkg/event"
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type Effect interface {
	core.RecordProxy
	ID() string
	Name() string
	Type() string
	Subscribe(EffectContext, EffectCallback) []event.Unsubscribe
	Verify(string) error
	DecodeValue(string) (any, error)
}

type EffectContext struct {
	User      User
	InvItemID string
}

type PersistentEffect interface {
	Subscribe(User) []event.Unsubscribe
}

var effectsList = map[string]EffectCreator{}
var persistentEffectsList = map[string]PersistentEffectCreator{}

type EffectCreator func(*core.Record) Effect
type PersistentEffectCreator func() PersistentEffect
type EffectCallback func()

func RegisterEffects(effects map[string]EffectCreator) {
	maps.Insert(effectsList, maps.All(effects))
}

func RegisterPersistentEffects(effects map[string]PersistentEffectCreator) {
	maps.Insert(persistentEffectsList, maps.All(effects))
}
