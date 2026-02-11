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

type TeleportToClosestCellByNameEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToClosestCellByNameEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if adventuria.GameActions.HasActionsInCategories(appCtx, ctx.User, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	canDone := adventuria.GameActions.CanDo(appCtx, ctx.User, "done")
	canDrop := adventuria.GameActions.CanDo(appCtx, ctx.User, "drop")

	if canDone && !canDrop {
		return false
	}

	return true
}

func (ef *TeleportToClosestCellByNameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*result.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				cellNames, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.MoveToClosestCellByNames(e.AppContext, cellNames...)
				if err != nil {
					return result.Err("internal error: failed to move to the cell by name"),
						fmt.Errorf("teleportToClosestCellByName: %w", err)
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToClosestCellByNameEffect) Verify(ctx adventuria.AppContext, value string) error {
	names, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToClosestCellByName: %w", err)
	}

	exp := make([]dbx.Expression, len(names))
	for i, name := range names {
		exp[i] = dbx.HashExp{"name": name}
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = ctx.App.RecordQuery(adventuria.GameCollections.Get(schema.CollectionCells)).
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

func (ef *TeleportToClosestCellByNameEffect) DecodeValue(value string) ([]string, error) {
	return strings.Split(value, ";"), nil
}

func (ef *TeleportToClosestCellByNameEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}
