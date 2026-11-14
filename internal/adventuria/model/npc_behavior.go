package model

type NPCBehaviorType string

type NPCBehaviorData struct {
	Id      string
	Type    NPCBehaviorType
	Payload []byte
	Data    []byte
}

type NPCBehavior struct {
	data NPCBehaviorData
}

func RestoreNPCBehavior(data NPCBehaviorData) *NPCBehavior {
	return &NPCBehavior{data: data}
}

func (b *NPCBehavior) ID() string {
	return b.data.Id
}

func (b *NPCBehavior) Type() NPCBehaviorType {
	return b.data.Type
}

func (b *NPCBehavior) Payload() []byte {
	return b.data.Payload
}

func (b *NPCBehavior) Data() []byte {
	return b.data.Data
}

func (b *NPCBehavior) SetData(data []byte) {
	b.data.Data = data
}
