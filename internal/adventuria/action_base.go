package adventuria

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ActionBase struct {
	core.BaseRecordProxy
	locator ServiceLocator
	user    User
}

func NewActionFromType(locator ServiceLocator, user User, actionType string) (Action, error) {
	actionCreator, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionCreator()

	actionsCollection, err := locator.Collections().Get(TableActions)
	if err != nil {
		return nil, err
	}

	action.SetProxyRecord(core.NewRecord(actionsCollection))
	action.setUser(user)
	action.setLocator(locator)
	action.SetType(actionType)

	return action, nil
}

func NewActionFromRecord(locator ServiceLocator, record *core.Record) Action {
	a := &ActionBase{
		locator: locator,
	}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionBase) CanDo() bool {
	panic("implement me")
}

func (a *ActionBase) Do(_ ActionRequest) (*ActionResult, error) {
	panic("implement me")
}

func (a *ActionBase) setUser(user User) {
	a.user = user
	a.Set("user", user.ID())
}

func (a *ActionBase) setLocator(locator ServiceLocator) {
	a.locator = locator
}

func (a *ActionBase) Save() error {
	return a.locator.PocketBase().Save(a)
}

func (a *ActionBase) SetType(t string) {
	a.Set("type", t)
}

func (a *ActionBase) User() User {
	return a.user
}

func (a *ActionBase) Locator() ServiceLocator {
	return a.locator
}

func (a *ActionBase) UserId() string {
	return a.GetString("user")
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

func NewLastUserAction(locator ServiceLocator, userId string) (Action, error) {
	a := &UserAction{
		ActionBase: ActionBase{locator: locator},
	}

	err := a.fetchLastUserAction(userId)
	if err != nil {
		return nil, err
	}
	a.bindHooks()

	return a, nil
}

func (ua *UserAction) bindHooks() {
	ua.locator.PocketBase().OnRecordAfterCreateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		notAffectNextStep := e.Record.GetBool("notAffectNextStep")
		if userId == ua.UserId() && !notAffectNextStep {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	ua.locator.PocketBase().OnRecordAfterUpdateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == ua.Id {
			ua.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	ua.locator.PocketBase().OnRecordAfterDeleteSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == ua.UserId() {
			ua.fetchLastUserAction(userId)
		}
		return e.Next()
	})
}

func (ua *UserAction) fetchLastUserAction(userId string) error {
	actions, err := ua.locator.PocketBase().FindRecordsByFilter(
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
		actionsCollection, _ := ua.locator.Collections().Get(TableActions)
		ua.SetProxyRecord(core.NewRecord(actionsCollection))
	}

	return nil
}
