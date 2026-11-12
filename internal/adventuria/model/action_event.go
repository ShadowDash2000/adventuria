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

func (c *ActionEventInfo) ID() string {
	return c.data.Id
}

func (c *ActionEventInfo) Name() string {
	return c.data.Name
}

func (c *ActionEventInfo) Type() ActionEventType {
	return c.data.Type
}

func (c *ActionEventInfo) Value() string {
	return c.data.Value
}
