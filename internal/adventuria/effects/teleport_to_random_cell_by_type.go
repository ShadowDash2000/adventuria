package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
	"fmt"
	"slices"
	"strings"
)

type TeleportToRandomCellByTypeEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToRandomCellByTypeEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *TeleportToRandomCellByTypeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	decodedValue, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return nil, err
	}

	switch decodedValue.Event {
	case "onAfterItemSave":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
				if e.Item.IDInventory() == ctx.InvItemID {
					err = teleportToRandomCellByType(ctx.User, decodedValue.CellTypes)
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
	case "onAfterItemUse":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
				if e.InvItemId == ctx.InvItemID {
					err = teleportToRandomCellByType(ctx.User, decodedValue.CellTypes)
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
	default:
		return nil, nil
	}
}

func teleportToRandomCellByType(user adventuria.User, cellTypes []string) error {
	_, err := user.MoveToClosestCellType(adventuria.CellType(helper.RandomItemFromSlice(cellTypes)))
	return err
}

func (ef *TeleportToRandomCellByTypeEffect) Verify(value string) error {
	events := []string{
		"onAfterItemSave",
		"onAfterItemUse",
	}

	decodedValue, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("teleportToRandomCellByType: %w", err)
	}

	if ok := slices.Contains(events, decodedValue.Event); !ok {
		return fmt.Errorf("teleportToRandomCellByType: unknown event: %s", decodedValue.Event)
	}

	for _, cellType := range decodedValue.CellTypes {
		if !adventuria.IsCellTypeExist(adventuria.CellType(cellType)) {
			return fmt.Errorf("teleportToRandomCellByType: unknown cell type: %s", cellType)
		}
	}

	return nil
}

type TeleportToRandomCellByTypeEffectValue struct {
	CellTypes []string
	Event     string
}

func (ef *TeleportToRandomCellByTypeEffect) DecodeValue(value string) (*TeleportToRandomCellByTypeEffectValue, error) {
	values := strings.Split(value, ";")

	if len(values) < 2 {
		return nil, fmt.Errorf("teleportToRandomCellByType: invalid value, expected format 'cellType1;cellType2;...;event'")
	}

	return &TeleportToRandomCellByTypeEffectValue{
		CellTypes: values[:len(values)-1],
		Event:     values[len(values)-1],
	}, nil
}

func (ef *TeleportToRandomCellByTypeEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}
