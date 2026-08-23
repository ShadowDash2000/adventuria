package items

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
	GetByFilter(ctx context.Context, filter Filter) ([]*model.Item, error)
	GetIDsByFilter(ctx context.Context, filter Filter) ([]string, error)
	GetByName(ctx context.Context, name string) (*model.Item, error)
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

func (i *Items) GetByFilter(ctx context.Context, filter Filter) ([]*model.Item, error) {
	return i.repository.GetByFilter(ctx, filter)
}

func (i *Items) GetIDsByFilter(ctx context.Context, filter Filter) ([]string, error) {
	return i.repository.GetIDsByFilter(ctx, filter)
}

func (i *Items) GetByName(ctx context.Context, name string) (*model.Item, error) {
	return i.repository.GetByName(ctx, name)
}

func (i *Items) GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error) {
	return i.GetIDsByFilter(ctx, Filter{
		ItemType:         t,
		Disabled:         new(false),
		IsRollable:       new(true),
		PriceGreaterThan: new(0),
	})
}

func (i *Items) GetAllRollable(ctx context.Context) ([]*model.Item, error) {
	return i.GetByFilter(ctx, Filter{
		Disabled:   new(false),
		IsRollable: new(true),
	})
}

func (i *Items) GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error) {
	return i.GetByFilter(ctx, Filter{
		ItemType:   t,
		Disabled:   new(false),
		IsRollable: new(true),
	})
}
