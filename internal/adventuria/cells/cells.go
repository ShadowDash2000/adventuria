package cells

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.CellInfo, error)
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
	GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error)
	GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
	CountGlobal(ctx context.Context) (int, error)
	GetAllByWorldID(ctx context.Context, worldId string) ([]*model.CellInfo, error)
}

type Cells struct {
	repository repository
}

func NewCells(repository repository) *Cells {
	return &Cells{
		repository: repository,
	}
}

func toCell(cell *model.CellInfo) (model.Cell, error) {
	cellDef, ok := Get(cell.Type())
	if !ok {
		return nil, fmt.Errorf("%w: %s", errs.ErrUnknownCellType, cell.Type())
	}
	return cellDef.new(*cell), nil
}

func (c *Cells) GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
	return c.GetByID(ctx, player.LastAction().Cell())
}

func (c *Cells) GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error) {
	return c.GetByIDWrapped(ctx, player.LastAction().Cell())
}

func (c *Cells) GetByID(ctx context.Context, id string) (*model.CellInfo, error) {
	return c.repository.GetByID(ctx, id)
}

func (c *Cells) GetByIDWrapped(ctx context.Context, id string) (model.Cell, error) {
	cell, err := c.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCell(cell)
}

func (c *Cells) GetByIDs(ctx context.Context, ids []string) ([]*model.CellInfo, error) {
	return c.repository.GetByIDs(ctx, ids)
}

func (c *Cells) GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error) {
	return c.repository.GetByLocalOrder(ctx, worldId, order)
}

func (c *Cells) GetByLocalOrderWrapped(ctx context.Context, worldId string, order int) (model.Cell, error) {
	cell, err := c.GetByLocalOrder(ctx, worldId, order)
	if err != nil {
		return nil, err
	}
	return toCell(cell)
}

func (c *Cells) GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error) {
	return c.repository.GetByGlobalOrder(ctx, order)
}

func (c *Cells) GetByGlobalOrderWrapped(ctx context.Context, order int) (model.Cell, error) {
	cell, err := c.repository.GetByGlobalOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return toCell(cell)
}

func (c *Cells) CountLocal(ctx context.Context, worldId string) (int, error) {
	return c.repository.CountLocal(ctx, worldId)
}

func (c *Cells) CountGlobal(ctx context.Context) (int, error) {
	return c.repository.CountGlobal(ctx)
}

func (c *Cells) GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error) {
	return c.repository.GetAllGlobalByType(ctx, t)
}

func (c *Cells) GetAllByWorldID(ctx context.Context, worldId string) ([]*model.CellInfo, error) {
	return c.repository.GetAllByWorldID(ctx, worldId)
}
