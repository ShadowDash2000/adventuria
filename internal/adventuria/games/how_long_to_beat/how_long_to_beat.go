package how_long_to_beat

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/levenshtein"
	"adventuria/pkg/mathhelper"
	"context"
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type repository interface {
	Create(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error)
	Update(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error)
	ExistsByIdDb(ctx context.Context, idDb int) (bool, error)
	GetByNamePartsAndYear(ctx context.Context, nameParts []string, year int) ([]*model.HowLongToBeat, error)
}

type remoteRepository interface {
	FetchLatestRelease(ctx context.Context) ([]*HowLongToBeatResponse, error)
}

type HowLongToBeat struct {
	repository       repository
	remoteRepository remoteRepository
}

func NewHowLongToBeat(repository repository, remoteRepository remoteRepository) *HowLongToBeat {
	return &HowLongToBeat{
		repository:       repository,
		remoteRepository: remoteRepository,
	}
}

func (h *HowLongToBeat) Parse(ctx context.Context) error {
	games, err := h.remoteRepository.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, game := range games {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ok, err := h.repository.ExistsByIdDb(ctx, game.ID)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		hltb, err := model.NewHowLongToBeat(model.HowLongToBeatCreate{
			IdDb:     game.ID,
			Name:     game.Name,
			Year:     game.ReleaseWorld,
			Campaign: math.Round(float64(game.CompMain) / 3600),
		})
		if err != nil {
			return err
		}

		_, err = h.Save(ctx, hltb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HowLongToBeat) Save(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	if hltb.IsNew() {
		return h.repository.Create(ctx, hltb)
	}

	return h.repository.Update(ctx, hltb)
}

func (h *HowLongToBeat) GetByNameAndYear(ctx context.Context, name string, year int) (*model.HowLongToBeat, error) {
	normalizedName := normalizeTitle(name)
	parts := strings.Fields(normalizedName)
	if len(parts) == 0 {
		return nil, errors.New("game name is empty")
	}

	hltbs, err := h.repository.GetByNamePartsAndYear(ctx, parts, year)
	if err != nil {
		return nil, err
	}

	type match struct {
		hltb      *model.HowLongToBeat
		exact     bool
		distance  int
		diffLen   int
		yearMatch bool
	}

	matches := make([]match, len(hltbs))

	for i, hltb := range hltbs {
		dbName := normalizeTitle(hltb.Name())
		dbYear := hltb.Year()

		matches[i] = match{
			hltb:      hltb,
			exact:     dbName == normalizedName,
			distance:  levenshtein.Distance(normalizedName, dbName),
			diffLen:   mathhelper.Abs(len(normalizedName) - len(dbName)),
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

	return matches[0].hltb, nil
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
