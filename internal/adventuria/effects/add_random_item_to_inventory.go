package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type AddRandomItemToInventoryEffect struct {
	adventuria.EffectRecord
}

func (ef *AddRandomItemToInventoryEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *AddRandomItemToInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) (*event.Result, error) {
			if e.ActionType == "buyItem" {
				idsAny, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.Inventory().AddItemById(helper.RandomItemFromSlice(idsAny.([]string)))
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "failed to add item to inventory",
					}, fmt.Errorf("addRandomItemToInventory: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *AddRandomItemToInventoryEffect) Verify(value string) error {
	decodedValue, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("addRandomItemToInventory: %w", err)
	}
	ids := decodedValue.([]string)

	idsAny := make([]any, len(ids))
	for i, id := range ids {
		idsAny[i] = id
	}

	var records []*core.Record
	err = adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionItems)).
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

func (ef *AddRandomItemToInventoryEffect) DecodeValue(value string) (any, error) {
	return strings.Split(value, ";"), nil
}
