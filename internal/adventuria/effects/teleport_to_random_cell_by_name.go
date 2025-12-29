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

type TeleportToRandomCellByNameEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToRandomCellByNameEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *TeleportToRandomCellByNameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				namesAny, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.MoveToCellName(helper.RandomItemFromSlice(namesAny.([]string)))
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't move to cell by name",
					}, fmt.Errorf("teleportToRandomCellByNameEffect: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToRandomCellByNameEffect) Verify(value string) error {
	decodedValue, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToRandomCellByName: %w", err)
	}
	names := decodedValue.([]string)

	namesAny := make([]any, len(names))
	for i, name := range names {
		namesAny[i] = name
	}

	var records []*core.Record
	err = adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionCells)).
		Where(dbx.In("name", namesAny...)).
		Select("id").
		All(&records)
	if err != nil {
		return fmt.Errorf("teleportToRandomCellByName: %w", err)
	}

	if len(names) != len(records) {
		return errors.New("teleportToRandomCellByName: not all cells found")
	}

	return nil
}

func (ef *TeleportToRandomCellByNameEffect) DecodeValue(value string) (any, error) {
	return strings.Split(value, ";"), nil
}
