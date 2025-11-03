package cells

import (
	"adventuria/internal/adventuria"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type CellShop struct {
	adventuria.CellBase
}

func NewCellShop() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellShop{
			adventuria.CellBase{},
		}
	}
}

func (c *CellShop) Roll(user adventuria.User) (*adventuria.WheelRollResult, error) {
	return nil, nil
}

func (c *CellShop) OnCellReached(user adventuria.User) error {
	var records []*core.Record
	err := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionItems)).
		Where(dbx.And(
			dbx.NewExp("type = \"buff\""),
			dbx.NewExp("isRollable = true"),
		)).
		Limit(6).
		All(&records)
	if err != nil {
		return fmt.Errorf("shop.OnCellReached: %w", err)
	}

	res := make([]string, len(records))
	for i, record := range records {
		res[i] = record.Id
	}

	user.LastAction().SetItemsList(res)
	user.LastAction().SetCanMove(true)

	return nil
}
