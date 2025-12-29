package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type TeleportToClosestCellByNameEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToClosestCellByNameEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *TeleportToClosestCellByNameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				namesAny, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.MoveToClosestCellByNames(namesAny.([]string)...)
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't move to cell by name",
					}, fmt.Errorf("teleportToClosestCellByName: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToClosestCellByNameEffect) Verify(value string) error {
	decodedValue, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToClosestCellByName: %w", err)
	}
	names := decodedValue.([]string)

	exp := make([]dbx.Expression, len(names))
	for i, name := range names {
		exp[i] = dbx.HashExp{"name": name}
	}

	var records []*core.Record
	err = adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionCells)).
		Where(dbx.Or(exp...)).
		Select("id").
		All(&records)
	if err != nil {
		return fmt.Errorf("teleportToClosestCellByName: %w", err)
	}

	if len(names) != len(records) {
		return errors.New("teleportToClosestCellByName: not all cells found")
	}

	return nil
}

func (ef *TeleportToClosestCellByNameEffect) DecodeValue(value string) (any, error) {
	return strings.Split(value, ";"), nil
}
