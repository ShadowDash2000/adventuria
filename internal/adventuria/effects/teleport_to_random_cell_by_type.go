package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"fmt"
	"strings"
)

type TeleportToRandomCellByTypeEffect struct {
	adventuria.EffectBase
}

func (ef *TeleportToRandomCellByTypeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				cellTypesAny, _ := ef.DecodeValue(ef.GetString("value"))
				_, err := ctx.User.MoveToClosestCellType(
					adventuria.CellType(helper.RandomItemFromSlice(cellTypesAny.([]string))),
				)
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "",
					}, fmt.Errorf("teleportToRandomCellByTypeEffect: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToRandomCellByTypeEffect) Verify(value string) error {
	decodedValue, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToRandomCellByType: %w", err)
	}

	for _, cellType := range decodedValue.([]string) {
		if !adventuria.IsCellTypeExist(adventuria.CellType(cellType)) {
			return fmt.Errorf("teleportToRandomCellByType: unknown cell type: %s", cellType)
		}
	}

	return nil
}

func (ef *TeleportToRandomCellByTypeEffect) DecodeValue(value string) (any, error) {
	return strings.Split(value, ";"), nil
}
