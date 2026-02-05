package adventuria

import "github.com/pocketbase/pocketbase/core"

type Inventories struct {
}

func NewInventories(ctx AppContext) *Inventories {
	i := &Inventories{}
	i.bindHooks(ctx)
	return i
}

func (i *Inventories) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")

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
}
