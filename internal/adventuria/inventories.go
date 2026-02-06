package adventuria

import (
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Inventories struct {
}

func NewInventories(ctx AppContext) *Inventories {
	i := &Inventories{}
	i.bindHooks(ctx)
	return i
}

func (i *Inventories) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString(schema.InventorySchema.User)

		ctx := AppContext{App: e.App}
		user, err := GameUsers.GetByID(ctx, userId)
		if err != nil {
			return e.Next()
		}

		if user.Inventory().HasItem(e.Record.Id) {
			return e.Next()
		}

		item, err := NewItemFromInventoryRecord(ctx, user, e.Record)
		if err != nil {
			return err
		}
		user.Inventory().RegisterItem(item)

		if _, err = user.OnAfterItemSave().Trigger(&OnAfterItemSave{
			AppContext: ctx,
			Item:       item,
		}); err != nil {
			e.App.Logger().Error(
				"Failed to trigger OnAfterItemSave event",
				"error", err,
			)
		}

		return e.Next()
	})
	ctx.App.OnRecordCreate(schema.CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetBool(schema.InventorySchema.IsActive) {
			e.Record.Set(schema.InventorySchema.Activated, types.NowDateTime())
		}
		return e.Next()
	})
}
