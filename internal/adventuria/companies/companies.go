package companies

import (
	"adventuria/internal/adventuria/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	GetOrCreate(ctx context.Context, id uuid.UUID, data model.CompanyCreate) (*model.Company, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, company *model.Company) (*model.Company, error)
}

type Companies struct {
	repository repository
}

func NewCompanies(repo repository) *Companies {
	return &Companies{
		repository: repo,
	}
}

func (c *Companies) GetOrCreate(ctx context.Context, id uuid.UUID, data model.CompanyCreate) (*model.Company, error) {
	return c.repository.GetOrCreate(ctx, id, data)
}

func (c *Companies) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return c.repository.GetChecksumsByIDs(ctx, ids)
}

func (c *Companies) Save(ctx context.Context, company *model.Company) (*model.Company, error) {
	return c.repository.Save(ctx, company)
}
