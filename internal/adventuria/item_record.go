package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

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
	return i.itemRecord.GetString(schema.ItemSchema.Name)
}

func (i *ItemRecordBase) Icon() string {
	return i.itemRecord.GetString(schema.ItemSchema.Icon)
}

func (i *ItemRecordBase) IsUsingSlot() bool {
	return i.itemRecord.GetBool(schema.ItemSchema.IsUsingSlot)
}

func (i *ItemRecordBase) IsActiveByDefault() bool {
	return i.itemRecord.GetBool(schema.ItemSchema.IsActiveByDefault)
}

func (i *ItemRecordBase) CanDrop() bool {
	return i.itemRecord.GetBool(schema.ItemSchema.CanDrop)
}

func (i *ItemRecordBase) IsRollable() bool {
	return i.itemRecord.GetBool(schema.ItemSchema.IsRollable)
}

func (i *ItemRecordBase) Order() int {
	return i.itemRecord.GetInt(schema.ItemSchema.Order)
}

func (i *ItemRecordBase) Type() string {
	return i.itemRecord.GetString(schema.ItemSchema.Type)
}

func (i *ItemRecordBase) Price() int {
	return i.itemRecord.GetInt(schema.ItemSchema.Price)
}
