package cells

import (
	"adventuria/internal/adventuria"
	"fmt"
	"math/rand/v2"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const shopMaxItems = 6

type CellShop struct {
	adventuria.CellRecord
}

func (c *CellShop) OnCellReached(ctx *adventuria.CellReachedContext) error {
	var records []*core.Record
	err := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionItems)).
		Where(dbx.And(
			dbx.NewExp("type = \"buff\""),
			dbx.NewExp("isRollable = true"),
		)).
		All(&records)
	if err != nil {
		return fmt.Errorf("shop.OnCellReached: %w", err)
	}

	if len(records) == 0 {
		return nil
	}

	res := make([]string, shopMaxItems)
	for i := 0; i < shopMaxItems; i++ {
		res[i] = records[rand.N(len(records))].Id
	}

	ctx.User.LastAction().SetItemsList(res)
	ctx.User.LastAction().SetCanMove(true)

	return nil
}

func (c *CellShop) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellShop) Verify(_ string) error {
	return nil
}
