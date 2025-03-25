package adventuria

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

const (
	ActionTypeRoll            = "roll"
	ActionTypeReroll          = "reroll"
	ActionTypeDrop            = "drop"
	ActionTypeChooseResult    = "chooseResult"
	ActionTypeChooseGame      = "chooseGame"
	ActionTypeRollCell        = "rollCell"
	ActionTypeRollItem        = "rollItem"
	ActionTypeRollWheelPreset = "rollWheelPreset"
)

type Action interface {
	core.RecordProxy
	Save() error
	UserId() string
	CellId() string
	SetCell(cellId string)
	Comment() string
	SetComment(comment string)
	Value() string
	SetValue(value any)
	Type() string
	SetIcon(*filesystem.File)
	SetNotAffectNextStep(bool)
}

type BaseAction struct {
	core.BaseRecordProxy
	gc *GameComponents
}

func NewAction(userId string, actionType string, gc *GameComponents) Action {
	a := &BaseAction{gc: gc}

	actionsCollection, _ := gc.cols.Get(TableActions)

	a.SetProxyRecord(core.NewRecord(actionsCollection))
	a.setUserId(userId)
	a.setType(actionType)

	return a
}

func (a *BaseAction) Save() error {
	return a.gc.app.Save(a)
}

func (a *BaseAction) UserId() string {
	return a.GetString("user")
}

func (a *BaseAction) setUserId(userId string) {
	a.Set("user", userId)
}

func (a *BaseAction) CellId() string {
	return a.GetString("cell")
}

func (a *BaseAction) SetCell(cellId string) {
	a.Set("cell", cellId)
}

func (a *BaseAction) Comment() string {
	return a.GetString("comment")
}

func (a *BaseAction) SetComment(comment string) {
	a.Set("comment", comment)
}

func (a *BaseAction) Value() string {
	return a.GetString("value")
}

func (a *BaseAction) SetValue(value any) {
	a.Set("value", value)
}

func (a *BaseAction) Type() string {
	return a.GetString("type")
}

func (a *BaseAction) setType(t string) {
	a.Set("type", t)
}

func (a *BaseAction) SetIcon(icon *filesystem.File) {
	if icon != nil {
		a.Set("icon", icon)
	}
}

func (a *BaseAction) SetNotAffectNextStep(b bool) {
	a.Set("notAffectNextStep", b)
}

type UserAction struct {
	BaseAction
}

func NewLastUserAction(userId string, gc *GameComponents) (Action, error) {
	a := &UserAction{
		BaseAction: BaseAction{gc: gc},
	}

	err := a.fetchLastUserAction(userId)
	if err != nil {
		return nil, err
	}
	a.bindHooks()

	return a, nil
}

func (ua *UserAction) bindHooks() {
	ua.gc.app.OnRecordAfterCreateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		notAffectNextStep := e.Record.GetBool("notAffectNextStep")
		if userId == ua.UserId() && !notAffectNextStep {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	ua.gc.app.OnRecordAfterUpdateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == ua.Id {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	ua.gc.app.OnRecordAfterDeleteSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == ua.UserId() {
			ua.fetchLastUserAction(userId)
		}
		return e.Next()
	})
}

func (ua *UserAction) fetchLastUserAction(userId string) error {
	actions, err := ua.gc.app.FindRecordsByFilter(
		TableActions,
		"user.id = {:userId} && notAffectNextStep = false",
		"-created",
		1,
		0,
		dbx.Params{"userId": userId},
	)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		ua.SetProxyRecord(actions[0])
	} else {
		actionsCollection, _ := ua.gc.cols.Get(TableActions)
		ua.SetProxyRecord(core.NewRecord(actionsCollection))
		ua.setUserId(userId)
	}

	return nil
}
