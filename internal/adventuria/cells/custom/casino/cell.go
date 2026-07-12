package casino

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ model.Refreshable = (*CellCasino)(nil)

const Type model.CellType = "casino"

type CellCasino struct {
	cells.CellBase
	items items
}

func NewDef(
	items items,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellCasino{
				CellBase: cells.NewCellBase(cell),
				items:    items,
			}
		},
		categories...,
	)
}

func (c *CellCasino) OnCellReached(_ context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	err := player.Progress().ItemWheelsCountChange(1)
	if err != nil {
		return err
	}

	player.Progress().SetCanMove(true)

	return c.refreshItems(player)
}

func (c *CellCasino) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

func (c *CellCasino) RefreshItems(_ context.Context, _ *model.Events, player *model.Player) error {
	return c.refreshItems(player)
}

func (c *CellCasino) refreshItems(player *model.Player) error {
	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}
	player.LastAction().SetItemsList(decodedValue.ItemIds)
	return nil
}
