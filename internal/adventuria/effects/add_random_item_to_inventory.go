package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"adventuria/pkg/result"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
)

type AddRandomItemToInventoryEffect struct {
	adventuria.EffectRecord
}

func (ef *AddRandomItemToInventoryEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *AddRandomItemToInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) (*result.Result, error) {
			if e.ActionType == "buyItem" {
				ids, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.Inventory().AddItemById(e.AppContext, helper.RandomItemFromSlice(ids))
				if err != nil {
					return result.Err("internal error: failed to add item to the inventory"),
						fmt.Errorf("addRandomItemToInventory: %w", err)
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *AddRandomItemToInventoryEffect) Verify(ctx adventuria.AppContext, value string) error {
	ids, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("addRandomItemToInventory: %w", err)
	}

	idsAny := make([]any, len(ids))
	for i, id := range ids {
		idsAny[i] = id
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = ctx.App.RecordQuery(adventuria.GameCollections.Get(schema.CollectionItems)).
		Where(dbx.In("id", idsAny...)).
		Select("id").
		All(&records)
	if err != nil {
		return fmt.Errorf("addRandomItemToInventory: %w", err)
	}

	if len(ids) != len(records) {
		return errors.New("addRandomItemToInventory: not all items found")
	}

	return nil
}

func (ef *AddRandomItemToInventoryEffect) DecodeValue(value string) ([]string, error) {
	return strings.Split(value, ";"), nil
}

func (ef *AddRandomItemToInventoryEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}
