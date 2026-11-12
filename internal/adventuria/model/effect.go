package model

type EffectType string

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
