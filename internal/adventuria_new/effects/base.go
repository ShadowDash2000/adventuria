package effects

import "adventuria/internal/adventuria_new/model"

type EffectBase struct {
	*model.EffectInfo
}

func NewEffectBase(info model.EffectInfo) EffectBase {
	return EffectBase{&info}
}

func (e EffectBase) Data() *model.EffectInfo {
	return e.EffectInfo
}
