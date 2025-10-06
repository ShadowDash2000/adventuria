package adventuria

import "github.com/pocketbase/pocketbase/core"

type ItemRecordBase struct {
	itemRecord core.BaseRecordProxy
}

func NewItemFromRecord(itemRecord *core.Record) ItemRecord {
	item := &ItemRecordBase{}
	item.itemRecord.SetProxyRecord(itemRecord)
	return item
}

func (i *ItemRecordBase) ID() string {
	return i.itemRecord.Id
}

func (i *ItemRecordBase) Name() string {
	return i.itemRecord.GetString("name")
}

func (i *ItemRecordBase) Icon() string {
	return i.itemRecord.GetString("icon")
}

func (i *ItemRecordBase) IsUsingSlot() bool {
	return i.itemRecord.GetBool("isUsingSlot")
}

func (i *ItemRecordBase) IsActiveByDefault() bool {
	return i.itemRecord.GetBool("isActiveByDefault")
}

func (i *ItemRecordBase) CanDrop() bool {
	return i.itemRecord.GetBool("canDrop")
}

func (i *ItemRecordBase) IsRollable() bool {
	return i.itemRecord.GetBool("isRollable")
}

func (i *ItemRecordBase) Order() int {
	return i.itemRecord.GetInt("order")
}
