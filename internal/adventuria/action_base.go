package adventuria

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ActionBase struct {
	core.BaseRecordProxy
}

func NewAction(userId string, actionType string) Action {
	a := &ActionBase{}

	actionsCollection, _ := GameCollections.Get(TableActions)

	a.SetProxyRecord(core.NewRecord(actionsCollection))
	a.setUserId(userId)
	a.SetType(actionType)

	return a
}

func NewActionFromRecord(record *core.Record) Action {
	a := &ActionBase{}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionBase) Save() error {
	return GameApp.Save(a)
}

func (a *ActionBase) UserId() string {
	return a.GetString("user")
}

func (a *ActionBase) setUserId(userId string) {
	a.Set("user", userId)
}

func (a *ActionBase) CellId() string {
	return a.GetString("cell")
}

func (a *ActionBase) SetCell(cellId string) {
	a.Set("cell", cellId)
}

func (a *ActionBase) Comment() string {
	return a.GetString("comment")
}

func (a *ActionBase) SetComment(comment string) {
	a.Set("comment", comment)
}

func (a *ActionBase) Value() string {
	return a.GetString("value")
}

func (a *ActionBase) SetValue(value any) {
	a.Set("value", value)
}

func (a *ActionBase) Type() string {
	return a.GetString("type")
}

func (a *ActionBase) SetType(t string) {
	a.Set("type", t)
}

func (a *ActionBase) SetNotAffectNextStep(b bool) {
	a.Set("notAffectNextStep", b)
}

func (a *ActionBase) CollectionRef() string {
	return a.GetString("collectionRef")
}

func (a *ActionBase) SetCollectionRef(collectionRef string) {
	a.Set("collectionRef", collectionRef)
}

func (a *ActionBase) DiceRoll() int {
	return a.GetInt("diceRoll")
}

func (a *ActionBase) SetDiceRoll(roll int) {
	a.Set("diceRoll", roll)
}

func (a *ActionBase) ItemsUsed() []string {
	return a.GetStringSlice("itemsUsed")
}

func (a *ActionBase) SetItemsUsed(items []string) {
	a.Set("itemsUsed", items)
}

type UserAction struct {
	ActionBase
}

func NewLastUserAction(userId string) (Action, error) {
	a := &UserAction{
		ActionBase: ActionBase{},
	}

	err := a.fetchLastUserAction(userId)
	if err != nil {
		return nil, err
	}
	a.bindHooks()

	return a, nil
}

func (ua *UserAction) bindHooks() {
	GameApp.OnRecordAfterCreateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		notAffectNextStep := e.Record.GetBool("notAffectNextStep")
		if userId == ua.UserId() && !notAffectNextStep {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	GameApp.OnRecordAfterUpdateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == ua.Id {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	GameApp.OnRecordAfterDeleteSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == ua.UserId() {
			ua.fetchLastUserAction(userId)
		}
		return e.Next()
	})
}

func (ua *UserAction) fetchLastUserAction(userId string) error {
	actions, err := GameApp.FindRecordsByFilter(
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
		actionsCollection, _ := GameCollections.Get(TableActions)
		ua.SetProxyRecord(core.NewRecord(actionsCollection))
		ua.setUserId(userId)
	}

	return nil
}
