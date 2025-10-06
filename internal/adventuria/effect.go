package adventuria

import (
	"adventuria/pkg/event"
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type Effect interface {
	core.RecordProxy
	setUser(User)
	ID() string
	Name() string
	Type() string
	User() User
	Subscribe(EffectCallback)
	PoolUnsubscribers(...event.Unsubscribe)
	Unsubscribe()
}

type PersistentEffect interface {
	User() User
	Subscribe()
	PoolUnsubscribers(...event.Unsubscribe)
	Unsubscribe()
}

var effectsList = map[string]EffectCreator{}
var persistentEffectsList = map[string]PersistentEffectCreator{}

type EffectCreator func(User, *core.Record) Effect
type PersistentEffectCreator func(User) PersistentEffect
type EffectCallback func()

func RegisterEffects(effects map[string]EffectCreator) {
	maps.Insert(effectsList, maps.All(effects))
}

func RegisterPersistentEffects(effects map[string]PersistentEffectCreator) {
	maps.Insert(persistentEffectsList, maps.All(effects))
}
