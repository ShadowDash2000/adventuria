package paid_movement_in_radius

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*PaidMovementInRadius)(nil)

type cellView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (p *PaidMovementInRadius) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	effectValue, err := p.decodeValue(p.Value())
	if err != nil {
		return nil, err
	}

	currentCell, err := p.cells.GetByID(ctx, player.LastAction().Cell())
	if err != nil {
		return nil, err
	}

	currentCellOrder := currentCell.LocalOrder()
	currentWorldId := currentCell.World()

	startOrder := currentCellOrder - effectValue.Radius
	if startOrder < 0 {
		startOrder = 0
	}

	worldCellsCount, err := p.cells.CountLocal(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	endOrder := currentCellOrder + effectValue.Radius
	if endOrder > worldCellsCount-1 {
		endOrder = worldCellsCount - 1
	}

	cells := make([]*model.CellInfo, endOrder-startOrder)
	j := 0
	for i := startOrder; i <= endOrder; i++ {
		if i == currentCellOrder {
			continue
		}

		cell, err := p.cells.GetByLocalOrder(ctx, currentWorldId, i)
		if err != nil {
			return nil, err
		}

		cells[j] = cell

		j++
	}

	return cellInfosToCellViews(cells), nil
}

func cellInfoToCellView(cell *model.CellInfo) *cellView {
	return &cellView{
		Id:   cell.ID(),
		Name: cell.Name(),
	}
}

func cellInfosToCellViews(cells []*model.CellInfo) []*cellView {
	cellsView := make([]*cellView, len(cells))
	for i, cell := range cells {
		cellsView[i] = cellInfoToCellView(cell)
	}
	return cellsView
}
