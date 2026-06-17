package model

import (
	"adventuria/pkg/event_new"
	"context"
)

type EffectType string
type EffectCallback func(ctx context.Context)

type Effect interface {
	Data() *EffectInfo
	CanUse(ctx context.Context, events *Events, player *Player) bool
	Subscribe(ctx context.Context, events *Events, player *Player, effectCtx EffectContext, callback EffectCallback) ([]event_new.Unsubscribe, error)
}

type EffectPersistent interface {
	Subscribe(ctx context.Context, events *Events, player *Player) ([]event_new.Unsubscribe, error)
}

type EffectContext struct {
	InvItemID string
	Priority  int
}

type EffectData struct {
	Id    string
	Name  string
	Type  EffectType
	Value string
}

type EffectInfo struct {
	data EffectData
}

func RestoreEffectInfo(data EffectData) *EffectInfo {
	return &EffectInfo{data: data}
}

func (e *EffectInfo) ID() string {
	return e.data.Id
}

func (e *EffectInfo) Name() string {
	return e.data.Name
}

func (e *EffectInfo) Type() EffectType {
	return e.data.Type
}

func (e *EffectInfo) Value() string {
	return e.data.Value
}
