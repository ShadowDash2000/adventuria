package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type AddRandomItemToInventoryEffect struct {
	adventuria.EffectBase
}

func (ef *AddRandomItemToInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			if e.ActionType == actions.ActionTypeBuyItem {
				idsAny, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.Inventory().AddItemById(helper.RandomItemFromSlice(idsAny.([]string)))
				if err != nil {
					return err
				}

				callback()
			}

			return e.Next()
		}),
	}
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
