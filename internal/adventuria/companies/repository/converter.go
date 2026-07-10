package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToCompany(record *core.Record) *model.Company {
	return model.RestoreCompany(model.CompanyData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.CompanySchema.IdDb),
		Name:     record.GetString(schema.CompanySchema.Name),
		Checksum: record.GetString(schema.CompanySchema.Checksum),
	})
}

func CompanyToRecord(company *model.Company, record *core.Record) {
	record.Id = company.ID()
	record.Set(schema.CompanySchema.IdDb, company.IdDb())
	record.Set(schema.CompanySchema.Name, company.Name())
	record.Set(schema.CompanySchema.Checksum, company.Checksum())
}
