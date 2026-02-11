package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
)

type AddItemsToInventoryEffect struct {
	adventuria.EffectRecord
}

func (ef *AddItemsToInventoryEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *AddItemsToInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if e.InvItemId != ctx.InvItemID {
				return e.Next()
			}

			ids, err := ef.DecodeValue(ef.GetString("value"))
			if err != nil {
				return result.Err("internal error: failed to decode effect value"),
					fmt.Errorf("addItemsToInventory: %w", err)
			}

			for _, id := range ids {
				item, ok := adventuria.GameItems.GetById(id)
				if !ok {
					return result.Err(fmt.Sprintf("item with id %s not found", id)), nil
				}

				_, err = ctx.User.Inventory().AddItem(e.AppContext, item)
				if err != nil {
					return result.Err("internal error: failed to add item to the inventory"),
						fmt.Errorf("addItemsToInventory: %w", err)
				}
			}

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *AddItemsToInventoryEffect) Verify(ctx adventuria.AppContext, value string) error {
	ids, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("addItemsToInventory: %w", err)
	}

	exp := make([]dbx.Expression, len(ids))
	for i, id := range ids {
		exp[i] = dbx.HashExp{"id": id}
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = ctx.App.RecordQuery(adventuria.GameCollections.Get(schema.CollectionItems)).
		Select("id").
		Where(dbx.Or(exp...)).
		All(&records)
	if err != nil {
		return fmt.Errorf("addItemsToInventory: %w", err)
	}

	if len(ids) != len(records) {
		return errors.New("addItemsToInventory: not all items found")
	}

	return nil
}

func (ef *AddItemsToInventoryEffect) DecodeValue(value string) ([]string, error) {
	return strings.Split(value, ";"), nil
}

func (ef *AddItemsToInventoryEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}
