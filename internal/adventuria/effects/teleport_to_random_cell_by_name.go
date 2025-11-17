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
	adventuria.EffectBase
}

func (ef *TeleportToRandomCellByNameEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			namesAny, _ := ef.DecodeValue(ef.GetString("value"))
			err := user.MoveToCellId(helper.RandomItemFromSlice(namesAny.([]string)))
			if err != nil {
				return err
			}

			callback()

			return e.Next()
		}),
	}
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
	return strings.Split(value, ","), nil
}
