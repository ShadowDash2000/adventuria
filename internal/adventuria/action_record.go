package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

type ActionRecordBase struct {
	core.BaseRecordProxy
	activityFilter CustomActivityFilter
}

func NewActionRecordFromRecord(record *core.Record) ActionRecord {
	a := &ActionRecordBase{}

	a.SetProxyRecord(record)

	return a
}

func (a *ActionRecordBase) SetProxyRecord(record *core.Record) {
	a.BaseRecordProxy.SetProxyRecord(record)
	if a.GetString(schema.ActionSchema.CustomActivityFilter) == "null" {
		a.activityFilter = CustomActivityFilter{}
	} else {
		if err := a.UnmarshalJSONField(schema.ActionSchema.CustomActivityFilter, &a.activityFilter); err != nil {
			PocketBase.Logger().Error("Failed to unmarshal custom_activity_filter", "err", err)
		}
	}
}

func (a *ActionRecordBase) ID() string {
	return a.Id
}

func (a *ActionRecordBase) User() string {
	return a.GetString(schema.ActionSchema.User)
}

func (a *ActionRecordBase) SetUser(id string) {
	a.Set(schema.ActionSchema.User, id)
}

func (a *ActionRecordBase) CellId() string {
	return a.GetString(schema.ActionSchema.Cell)
}

func (a *ActionRecordBase) setCell(cellId string) {
	a.Set(schema.ActionSchema.Cell, cellId)
}

func (a *ActionRecordBase) Comment() string {
	return a.GetString(schema.ActionSchema.Comment)
}

func (a *ActionRecordBase) SetComment(comment string) {
	a.Set(schema.ActionSchema.Comment, comment)
}

func (a *ActionRecordBase) Activity() string {
	return a.GetString(schema.ActionSchema.Activity)
}

func (a *ActionRecordBase) SetActivity(id string) {
	a.Set(schema.ActionSchema.Activity, id)
}

func (a *ActionRecordBase) Type() ActionType {
	return ActionType(a.GetString(schema.ActionSchema.Type))
}

func (a *ActionRecordBase) SetType(t ActionType) {
	a.Set(schema.ActionSchema.Type, string(t))
}

func (a *ActionRecordBase) DiceRoll() int {
	return a.GetInt(schema.ActionSchema.DiceRoll)
}

func (a *ActionRecordBase) SetDiceRoll(roll int) {
	a.Set(schema.ActionSchema.DiceRoll, roll)
}

func (a *ActionRecordBase) UsedItems() []string {
	return a.GetStringSlice(schema.ActionSchema.UsedItems)
}

func (a *ActionRecordBase) UsedItemAppend(itemId string) {
	a.Set(schema.ActionSchema.UsedItems, append(a.UsedItems(), itemId))
}

func (a *ActionRecordBase) SetUsedItems(items []string) {
	a.Set(schema.ActionSchema.UsedItems, items)
}

func (a *ActionRecordBase) ItemsList() ([]string, error) {
	var items []string
	return items, a.UnmarshalJSONField(schema.ActionSchema.ItemsList, &items)
}

func (a *ActionRecordBase) SetItemsList(items []string) {
	a.Set(schema.ActionSchema.ItemsList, items)
}

func (a *ActionRecordBase) CanMove() bool {
	return a.GetBool(schema.ActionSchema.CanMove)
}

func (a *ActionRecordBase) SetCanMove(b bool) {
	a.Set(schema.ActionSchema.CanMove, b)
}

func (a *ActionRecordBase) CustomActivityFilter() *CustomActivityFilter {
	return &a.activityFilter
}

func (a *ActionRecordBase) ClearCustomActivityFilter() {
	a.activityFilter = CustomActivityFilter{}
}

// MarkAsNew resets the action record to a new state
// Note: after calling this method, the record will be saved as a new record
func (a *ActionRecordBase) MarkAsNew() {
	a.ProxyRecord().MarkAsNew()
	a.ProxyRecord().Set(schema.ActionSchema.Id, "")
	a.SetComment("")
	a.SetActivity("")
	a.SetDiceRoll(0)
	a.SetUsedItems([]string{})
}
