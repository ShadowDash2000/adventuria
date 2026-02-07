package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
)

type CellCasino struct {
	adventuria.CellRecord
}

type cellCasinoValue struct {
	ItemIds         []string `json:"item_ids"`
	PriceMultiplier float32  `json:"price_multiplier"`
}

func (c *CellCasino) OnCellReached(ctx *adventuria.CellReachedContext) error {
	ctx.User.SetItemWheelsCount(ctx.User.ItemWheelsCount() + 1)
	ctx.User.LastAction().SetCanMove(true)

	var decodedValue cellCasinoValue
	if err := c.UnmarshalJSONField("value", &decodedValue); err != nil {
		return err
	}
	ctx.User.LastAction().SetItemsList(decodedValue.ItemIds)

	return nil
}

func (c *CellCasino) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellCasino) Verify(ctx adventuria.AppContext, value string) error {
	var decodedValue cellCasinoValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return fmt.Errorf("item.verify: invalid JSON: %w", err)
	}

	if len(decodedValue.ItemIds) == 0 {
		return fmt.Errorf("item.verify: invalid value: %s", value)
	}

	exp := make([]dbx.Expression, len(decodedValue.ItemIds))
	for i, id := range decodedValue.ItemIds {
		exp[i] = dbx.HashExp{"id": id}
	}

	var records []struct {
		Id string `db:"id"`
	}
	err := ctx.App.
		RecordQuery(schema.CollectionItems).
		Select("id").
		Where(dbx.Or(exp...)).
		All(&records)
	if err != nil {
		return fmt.Errorf("item.verify: %w", err)
	}

	if len(decodedValue.ItemIds) != len(records) {
		return fmt.Errorf("item.verify: some of ids not found: %s", value)
	}

	return nil
}
