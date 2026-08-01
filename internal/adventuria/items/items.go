package items

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
	GetAllRollable(ctx context.Context) ([]*model.Item, error)
	GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error)
	GetIDsByFilter(ctx context.Context, filter Filter) ([]string, error)
}

type Items struct {
	repository repository
}

func NewItems(repository repository) *Items {
	return &Items{repository: repository}
}

func (i *Items) GetByID(ctx context.Context, id string) (*model.Item, error) {
	return i.repository.GetByID(ctx, id)
}

func (i *Items) GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error) {
	return i.repository.GetByIDs(ctx, ids)
}

func (i *Items) GetAllRollable(ctx context.Context) ([]*model.Item, error) {
	return i.repository.GetAllRollable(ctx)
}

func (i *Items) GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error) {
	return i.repository.GetAllRollableByType(ctx, t)
}

func (i *Items) GetIDsByFilter(ctx context.Context, filter Filter) ([]string, error) {
	return i.repository.GetIDsByFilter(ctx, filter)
}

func (i *Items) GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error) {
	return i.GetIDsByFilter(ctx, Filter{
		ItemType:         t,
		IsRollable:       new(true),
		PriceGreaterThan: new(0),
	})
}
