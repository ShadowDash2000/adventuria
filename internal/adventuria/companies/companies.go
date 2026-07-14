package companies

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByIdDb(ctx context.Context, idDb string) (*model.Company, error)
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

func (c *Companies) GetOrCreate(ctx context.Context, data model.CompanyCreate) (*model.Company, error) {
	company, err := c.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrCompanyNotFound) {
			return model.NewCompany(data)
		}
		return nil, err
	}

	return company, nil
}

func (c *Companies) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return c.repository.GetChecksumsByIDs(ctx, ids)
}

func (c *Companies) Save(ctx context.Context, company *model.Company) (*model.Company, error) {
	return c.repository.Save(ctx, company)
}
