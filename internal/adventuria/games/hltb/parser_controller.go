package hltb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type ParserController struct {
	parser *Parser
}

func New() (*ParserController, error) {
	p, err := NewParser()
	if err != nil {
		return nil, err
	}

	return &ParserController{
		parser: p,
	}, nil
}

func (p *ParserController) Parse(ctx context.Context, limit int) {
	if err := p.parseTime(ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse time", "error", err)
		return
	}
}

func (p *ParserController) parseTime(ctx context.Context, limit int) error {
	for {
		gameRecords, err := p.getGamesWithoutTime(limit)
		if err != nil {
			return err
		}

		if len(gameRecords) == 0 {
			break
		}

		for _, game := range gameRecords {
			time, err := p.parseWalkthroughTime(ctx, game)
			if err != nil {
				if errors.Is(err, ErrGameNotFound) {
					adventuria.PocketBase.Logger().Info("parseTime(): Game not found", "game", game.Name())
					game.SetCampaignTime(0)
				} else {
					adventuria.PocketBase.Logger().Error("parseTime(): Failed to parse time",
						"error", err,
						"game", game.Name(),
					)
					continue
				}
			} else {
				game.SetHltbID(time.GameID)
				game.SetCampaignTime(time.Campaign)
			}

			err = adventuria.PocketBase.Save(game.ProxyRecord())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *ParserController) parseWalkthroughTime(ctx context.Context, game games.GameRecord) (*WalkthroughTime, error) {
	steamAppId := game.SteamAppId()

	var (
		gameTime *WalkthroughTime
		err      error
	)
	if steamAppId > 0 {
		gameTime, err = p.parser.ParseBySteamAppId(ctx, steamAppId)

		if err != nil {
			adventuria.PocketBase.Logger().Debug(
				"parseTime(): Failed to parse time by steam app id",
				"error", err,
				"steamAppId", steamAppId,
			)
		}
	}

	if steamAppId <= 0 || gameTime == nil {
		normalizedTitle := normalizeTitle(game.Name())

		adventuria.PocketBase.Logger().Debug(
			"parseTime(): Parsing time",
			"game", game.Name(),
			"normalizedTitle", normalizedTitle,
		)

		gameTime, err = p.parser.ParseTime(ctx, normalizedTitle)
	}

	return gameTime, err
}

func (p *ParserController) getGamesWithoutTime(limit int) ([]games.GameRecord, error) {
	records, err := adventuria.PocketBase.FindRecordsByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		"campaign_time = -1",
		"",
		limit,
		0,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res := make([]games.GameRecord, len(records))
	for i, record := range records {
		res[i] = games.NewGameFromRecord(record)
	}

	return res, nil
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
