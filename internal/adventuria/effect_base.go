package adventuria

import (
	"adventuria/pkg/event"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

type EffectBase struct {
	core.BaseRecordProxy
	user          User
	unsubscribers []event.Unsubscribe
}

func NewEffectFromRecord(user User, record *core.Record) (Effect, error) {
	t := record.GetString("type")

	effectCreator, ok := effectsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", t)
	}

	effect := effectCreator(user, record)

	return effect, nil
}

func NewEffect(e Effect) EffectCreator {
	return func(user User, record *core.Record) Effect {
		e.setUser(user)
		e.SetProxyRecord(record)
		return e
	}
}

func (e *EffectBase) setUser(user User) {
	e.user = user
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

func (e *EffectBase) User() User {
	return e.user
}

func (e *EffectBase) Subscribe(_ EffectCallback) {
	panic("implement me")
}

func (e *EffectBase) PoolUnsubscribers(u ...event.Unsubscribe) {
	e.unsubscribers = append(e.unsubscribers, u...)
}

func (e *EffectBase) Unsubscribe() {
	for _, u := range e.unsubscribers {
		u()
	}
}

type PersistentEffectBase struct {
	core.BaseRecordProxy
	user          User
	unsubscribers []event.Unsubscribe
}

func NewPersistentEffect() PersistentEffectCreator {
	return func(user User) PersistentEffect {
		e := &PersistentEffectBase{
			user: user,
		}
		return e
	}
}

func (e *PersistentEffectBase) User() User {
	return e.user
}

func (e *PersistentEffectBase) Subscribe() {
	panic("implement me")
}

func (e *PersistentEffectBase) PoolUnsubscribers(u ...event.Unsubscribe) {
	e.unsubscribers = append(e.unsubscribers, u...)
}

func (e *PersistentEffectBase) Unsubscribe() {
	for _, u := range e.unsubscribers {
		u()
	}
}
