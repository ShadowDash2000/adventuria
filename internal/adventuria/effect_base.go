package adventuria

import (
	"adventuria/pkg/event"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// ensure EffectBase implements Effect
var _ Effect = (*EffectBase)(nil)

type EffectBase struct {
	core.BaseRecordProxy
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

func (e *EffectBase) ID() string {
	return e.Id
}

func (e *EffectBase) Name() string {
	return e.GetString("name")
}

func (e *EffectBase) Type() string {
	return e.GetString("type")
}

func (e *EffectBase) Subscribe(_ EffectContext, _ EffectCallback) ([]event.Unsubscribe, error) {
	panic("implement me")
}

func (e *EffectBase) Verify(_ string) error {
	return errors.New("effect verifier is not implemented")
}

func (e *EffectBase) DecodeValue(_ string) (any, error) {
	return nil, errors.New("effect decode value is not implemented")
}

// ensure PersistentEffectBase implements PersistentEffect
var _ PersistentEffect = (*PersistentEffectBase)(nil)

type PersistentEffectBase struct {
	core.BaseRecordProxy
}

func NewPersistentEffect(e PersistentEffect) PersistentEffectCreator {
	return func() PersistentEffect {
		return e
	}
}

func (e *PersistentEffectBase) Subscribe(_ User) []event.Unsubscribe {
	panic("implement me")
}
