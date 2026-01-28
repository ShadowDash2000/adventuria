package adventuria

import (
	"adventuria/pkg/event"
	"fmt"
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type Effect interface {
	core.RecordProxy
	ID() string
	Name() string
	Type() string
	CanUse(EffectContext) bool
	Subscribe(EffectContext, EffectCallback) ([]event.Unsubscribe, error)
	Verify(string) error
	GetVariants(EffectContext) any
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

func NewEffectFromRecord(record *core.Record) (Effect, error) {
	t := record.GetString("type")

	effectCreator, ok := effectsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", t)
	}

	effect := effectCreator(record)

	return effect, nil
}

func NewEffect(e Effect) EffectCreator {
	return func(record *core.Record) Effect {
		e.SetProxyRecord(record)
		return e
	}
}

func NewPersistentEffect(e PersistentEffect) PersistentEffectCreator {
	return func() PersistentEffect {
		return e
	}
}
