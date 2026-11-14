package model

type ActionEventType string

type ActionEventData struct {
	Id    string
	Name  string
	Type  ActionEventType
	Value string
}

type ActionEventInfo struct {
	data ActionEventData
}

func RestoreActionEvent(data ActionEventData) *ActionEventInfo {
	return &ActionEventInfo{data: data}
}

func (a *ActionEventInfo) Clone() *ActionEventInfo {
	if a == nil {
		return nil
	}

	return &ActionEventInfo{data: ActionEventData{
		Id:    a.data.Id,
		Name:  a.data.Name,
		Type:  a.data.Type,
		Value: a.data.Value,
	}}
}

func (a *ActionEventInfo) ID() string {
	return a.data.Id
}

func (a *ActionEventInfo) Name() string {
	return a.data.Name
}

func (a *ActionEventInfo) Type() ActionEventType {
	return a.data.Type
}

func (a *ActionEventInfo) Value() string {
	return a.data.Value
}
