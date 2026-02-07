package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"fmt"
	"math/rand/v2"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

var _ adventuria.CellRefreshable = (*CellShop)(nil)

const shopMaxItems = 6

type CellShop struct {
	adventuria.CellRecord
}

func (c *CellShop) RefreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	return c.refreshItems(ctx, user)
}

func (c *CellShop) OnCellReached(ctx *adventuria.CellReachedContext) error {
	if err := c.refreshItems(ctx.AppContext, ctx.User); err != nil {
		return err
	}

	ctx.User.LastAction().SetCanMove(true)

	return nil
}

func (c *CellShop) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellShop) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (c *CellShop) refreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	var records []*core.Record
	err := ctx.App.RecordQuery(adventuria.GameCollections.Get(schema.CollectionItems)).
		Where(dbx.And(
			dbx.NewExp(fmt.Sprintf("%s = \"buff\"", schema.ItemSchema.Type)),
			dbx.NewExp(fmt.Sprintf("%s = true", schema.ItemSchema.IsRollable)),
			dbx.NewExp(fmt.Sprintf("%s > 0", schema.ItemSchema.Price)),
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

	user.LastAction().SetItemsList(res)

	return nil
}
