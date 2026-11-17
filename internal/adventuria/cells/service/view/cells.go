package view

import "context"

type repository interface {
	GetAllBoard(ctx context.Context) ([]*CellBoard, error)
	GetDetailedByID(ctx context.Context, id string) (*CellDetailed, error)
}

type Cells struct {
	repository repository
}

func NewCells(repository repository) *Cells {
	return &Cells{
		repository: repository,
	}
}

func (c *Cells) GetAllBoard(ctx context.Context) ([]*CellBoard, error) {
	return c.repository.GetAllBoard(ctx)
}

func (c *Cells) GetDetailedByID(ctx context.Context, id string) (*CellDetailed, error) {
	return c.repository.GetDetailedByID(ctx, id)
}
