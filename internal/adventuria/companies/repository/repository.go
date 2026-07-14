package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Company, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionCompanies).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.CompanySchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCompanyNotFound
		}

		return nil, err
	}

	return RecordToCompany(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionCompanies).
		WithContext(ctx).
		Select(
			schema.CompanySchema.Id,
			schema.CompanySchema.Checksum,
		).
		Where(dbx.In(
			schema.CompanySchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.CompanySchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, company *model.Company) (*model.Company, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionCompanies)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	CompanyToRecord(company, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToCompany(record), nil
}

func (r *Repository) Update(ctx context.Context, company *model.Company) (*model.Company, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionCompanies, company.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	CompanyToRecord(company, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToCompany(record), nil
}

func (r *Repository) Save(ctx context.Context, company *model.Company) (*model.Company, error) {
	if company.IsNew() {
		return r.Create(ctx, company)
	}

	return r.Update(ctx, company)
}
