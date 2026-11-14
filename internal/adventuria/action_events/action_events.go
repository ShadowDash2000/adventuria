package action_events

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"fmt"
)

type ActionEvents struct{}

func NewActionEvents() *ActionEvents {
	return &ActionEvents{}
}

func (a *ActionEvents) CreateFromActionEventInfo(info *model.ActionEventInfo) (ActionEvent, error) {
	actionEventDef, ok := Get(info.Type())
	if !ok {
		return nil, fmt.Errorf("%w: %s", errs.ErrUnknownActionEventType, info.Type())
	}
	return actionEventDef.New(info)
}
