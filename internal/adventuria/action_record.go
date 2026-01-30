package adventuria

import "github.com/pocketbase/pocketbase/core"

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
	if a.GetString("custom_activity_filter") == "null" {
		a.activityFilter = CustomActivityFilter{}
	} else {
		if err := a.UnmarshalJSONField("custom_activity_filter", &a.activityFilter); err != nil {
			PocketBase.Logger().Error("Failed to unmarshal custom_activity_filter", "err", err)
		}
	}
}

func (a *ActionRecordBase) ID() string {
	return a.Id
}

func (a *ActionRecordBase) User() string {
	return a.GetString("user")
}

func (a *ActionRecordBase) SetUser(id string) {
	a.Set("user", id)
}

func (a *ActionRecordBase) CellId() string {
	return a.GetString("cell")
}

func (a *ActionRecordBase) setCell(cellId string) {
	a.Set("cell", cellId)
}

func (a *ActionRecordBase) Comment() string {
	return a.GetString("comment")
}

func (a *ActionRecordBase) SetComment(comment string) {
	a.Set("comment", comment)
}

func (a *ActionRecordBase) Activity() string {
	return a.GetString("activity")
}

func (a *ActionRecordBase) SetActivity(id string) {
	a.Set("activity", id)
}

func (a *ActionRecordBase) Type() ActionType {
	return ActionType(a.GetString("type"))
}

func (a *ActionRecordBase) SetType(t ActionType) {
	a.Set("type", string(t))
}

func (a *ActionRecordBase) DiceRoll() int {
	return a.GetInt("diceRoll")
}

func (a *ActionRecordBase) SetDiceRoll(roll int) {
	a.Set("diceRoll", roll)
}

func (a *ActionRecordBase) UsedItems() []string {
	return a.GetStringSlice("used_items")
}

func (a *ActionRecordBase) UsedItemAppend(itemId string) {
	a.Set("used_items", append(a.UsedItems(), itemId))
}

func (a *ActionRecordBase) SetUsedItems(items []string) {
	a.Set("used_items", items)
}

func (a *ActionRecordBase) ItemsList() ([]string, error) {
	var items []string
	return items, a.UnmarshalJSONField("items_list", &items)
}

func (a *ActionRecordBase) SetItemsList(items []string) {
	a.Set("items_list", items)
}

func (a *ActionRecordBase) CanMove() bool {
	return a.GetBool("can_move")
}

func (a *ActionRecordBase) SetCanMove(b bool) {
	a.Set("can_move", b)
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
	a.ProxyRecord().Set("id", "")
	a.SetComment("")
	a.SetActivity("")
	a.SetDiceRoll(0)
	a.SetUsedItems([]string{})
}
