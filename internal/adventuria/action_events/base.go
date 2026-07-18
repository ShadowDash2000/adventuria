package action_events

import "adventuria/internal/adventuria/model"

type ActionEventBase struct {
	*model.ActionEventInfo
}

func NewActionEventBase(actionEventInfo model.ActionEventInfo) ActionEventBase {
	return ActionEventBase{&actionEventInfo}
}

func (c ActionEventBase) Data() *model.ActionEventInfo {
	return c.ActionEventInfo
}
