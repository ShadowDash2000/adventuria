package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionHowLongToBeat)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	HowLongToBeatToRecord(hltb, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToHowLongToBeat(record), nil
}

func (r *Repository) Update(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionHowLongToBeat, hltb.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrHowLongToBeatNotFound
		}
		return nil, err
	}

	HowLongToBeatToRecord(hltb, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToHowLongToBeat(record), nil
}

func (r *Repository) Save(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	if hltb.IsNew() {
		return r.Create(ctx, hltb)
	}

	return r.Update(ctx, hltb)
}

func (r *Repository) ExistsByIdDb(ctx context.Context, idDb int) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionHowLongToBeat).
		WithContext(ctx).
		Select(schema.HowLongToBeatSchema.Id).
		Where(dbx.HashExp{schema.HowLongToBeatSchema.IdDb: idDb}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repository) GetByNameAndYear(ctx context.Context, name string, year int) (*model.HowLongToBeat, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	normalizedName := normalizeTitle(name)
	parts := strings.Fields(normalizedName)
	if len(parts) == 0 {
		return nil, errors.New("game name is empty")
	}

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionHowLongToBeat).
		WithContext(ctx).
		Where(dbx.Like(schema.HowLongToBeatSchema.Name, parts...)).
		AndWhere(dbx.Or(
			dbx.HashExp{schema.HowLongToBeatSchema.Year: year},
			dbx.HashExp{schema.HowLongToBeatSchema.Year: 0},
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errs.ErrHowLongToBeatNotFound
	}

	type match struct {
		record    *core.Record
		exact     bool
		distance  int
		diffLen   int
		yearMatch bool
	}

	matches := make([]match, len(records))

	for i, r := range records {
		dbName := normalizeTitle(r.GetString(schema.HowLongToBeatSchema.Name))
		dbYear := r.GetInt(schema.HowLongToBeatSchema.Year)

		matches[i] = match{
			record:    r,
			exact:     dbName == normalizedName,
			distance:  levenshteinDistance(normalizedName, dbName),
			diffLen:   int(math.Abs(float64(len(normalizedName) - len(dbName)))),
			yearMatch: dbYear == year,
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].exact != matches[j].exact {
			return matches[i].exact && !matches[j].exact
		}

		if matches[i].distance != matches[j].distance {
			return matches[i].distance < matches[j].distance
		}

		if matches[i].yearMatch != matches[j].yearMatch {
			return matches[i].yearMatch && !matches[j].yearMatch
		}

		return matches[i].diffLen < matches[j].diffLen
	})

	return RecordToHowLongToBeat(matches[0].record), nil
}

var (
	regParens = regexp.MustCompile(`\s*[(\[{].*?[)\]}]\s*`)
	regSpaces = regexp.MustCompile(`\s+`)
)

func normalizeTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	s = regParens.ReplaceAllString(s, " ")

	s = strings.ToLower(s)

	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			return r
		case r == '\'' || r == '.' || r == '/' || r == '\\':
			return r
		default:
			return ' '
		}
	}, s)

	s = regSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func levenshteinDistance(s1, s2 string) int {
	r1 := []rune(strings.ToLower(s1))
	r2 := []rune(strings.ToLower(s2))
	n, m := len(r1), len(r2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= m; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(matrix[i-1][j]+1, min(matrix[i][j-1]+1, matrix[i-1][j-1]+cost))
		}
	}
	return matrix[n][m]
}
