package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	seasonsRepo "adventuria/internal/adventuria/seasons/repository"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Seasons struct {
	pb core.App
}

func NewSeasons(pb core.App) *Seasons {
	return &Seasons{pb: pb}
}

func (s *Seasons) GetBySlug(ctx context.Context, slug string) (*model.Season, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, s.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionSeasons).
		WithContext(ctx).
		Where(dbx.HashExp{schema.SeasonSchema.Slug: slug}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSeasonNotFound
		}
		return nil, err
	}

	return seasonsRepo.RecordToSeason(&record), nil
}
