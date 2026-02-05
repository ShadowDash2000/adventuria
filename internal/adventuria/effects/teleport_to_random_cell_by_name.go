package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
)

type TeleportToRandomCellByNameEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToRandomCellByNameEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
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

func (ef *TeleportToRandomCellByNameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				cellNames, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.MoveToCellName(e.AppContext, helper.RandomItemFromSlice(cellNames))
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't move to cell by name",
					}, fmt.Errorf("teleportToRandomCellByNameEffect: %w", err)
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToRandomCellByNameEffect) Verify(ctx adventuria.AppContext, value string) error {
	names, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToRandomCellByName: %w", err)
	}

	namesAny := make([]any, len(names))
	for i, name := range names {
		namesAny[i] = name
	}

	var records []struct {
		Id string `db:"id"`
	}
	err = ctx.App.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionCells)).
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

func (ef *TeleportToRandomCellByNameEffect) DecodeValue(value string) ([]string, error) {
	return strings.Split(value, ";"), nil
}

func (ef *TeleportToRandomCellByNameEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}
