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

func (e *EffectInfo) Clone() *EffectInfo {
	if e == nil {
		return nil
	}

	return &EffectInfo{data: EffectData{
		Id:    e.data.Id,
		Name:  e.data.Name,
		Type:  e.data.Type,
		Value: e.data.Value,
	}}
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
