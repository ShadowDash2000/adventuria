package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
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
	err := ctx.App.RecordQuery(adventuria.GameCollections.Get(schema.CollectionItems)).
		Where(dbx.And(
			dbx.NewExp("type = \"buff\""),
			dbx.NewExp("isRollable = true"),
			dbx.NewExp("price > 0"),
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

	return ctx.App.Save(ctx.User.LastAction().ProxyRecord())
}

func (c *CellShop) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellShop) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}
