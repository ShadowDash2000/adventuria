package hltb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser *Parser
}

func New(r time.Duration, b int) (*ParserController, error) {
	p, err := NewParser(r, b)
	if err != nil {
		return nil, err
	}

	return &ParserController{
		parser: p,
	}, nil
}

func (p *ParserController) Parse(ctx context.Context, limit int64) {
	if err := p.parseTime(ctx, limit); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse time", "error", err)
		return
	}
}

func (p *ParserController) ParseWithWorkers(ctx context.Context, limit int64, workers int, waitEvery, wait time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_ = p.parser.RefreshToken(ctx, false)
	p.runRefreshToken(ctx)
	p.runParseTime(ctx, limit, workers, waitEvery, wait)
}

func (p *ParserController) runRefreshToken(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := p.parser.RefreshToken(ctx, false)
				if err != nil {
					adventuria.PocketBase.Logger().Error("Failed to refresh token", "error", err)
				}
			}
		}
	}()
}

// runParseTime runs the parseTimeFromChan method in goroutines.
// Method is blocking and will wait until all games are processed.
func (p *ParserController) runParseTime(ctx context.Context, limit int64, workers int, waitEvery, wait time.Duration) {
	ch := p.fetchGamesWithoutTime(ctx, limit, waitEvery, wait)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			p.parseTimeFromChan(ctx, ch)
		}()
	}

	wg.Wait()
}

// fetchGamesWithoutTime returns a channel of games without campaign time
// that will work until there are no more games to fetch.
func (p *ParserController) fetchGamesWithoutTime(ctx context.Context, limit int64, waitEvery, wait time.Duration) <-chan games.GameRecord {
	ch := make(chan games.GameRecord, limit)

	go func() {
		ticker := time.NewTicker(waitEvery)
		defer ticker.Stop()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			gameRecords, err := p.getGamesWithoutTime(ctx, limit)
			if err != nil {
				return
			}

			if len(gameRecords) == 0 {
				return
			}

			for _, game := range gameRecords {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					time.Sleep(wait)
				case ch <- game:
				}
			}
		}
	}()

	return ch
}

func (p *ParserController) parseTimeFromChan(ctx context.Context, ch <-chan games.GameRecord) {
	for game := range ch {
		time, err := p.parseWalkthroughTime(ctx, game)
		if err != nil {
			if errors.Is(err, ErrGameNotFound) {
				adventuria.PocketBase.Logger().Debug("parseTimeFromChan(): Game not found", "game", game.Name())
				game.SetCampaignTime(0)
			} else {
				adventuria.PocketBase.Logger().Error("parseTimeFromChan(): Failed to parse time",
					"error", err,
					"game", game.Name(),
				)
			}
		} else {
			game.SetHltbID(time.GameID)
			game.SetCampaignTime(time.Campaign)
		}

		err = adventuria.PocketBase.Save(game.ProxyRecord())
		if err != nil {
			adventuria.PocketBase.Logger().Error("parseTimeFromChan(): Failed to save game", "error", err)
		}
	}
}

func (p *ParserController) parseTime(ctx context.Context, limit int64) error {
	for {
		gameRecords, err := p.getGamesWithoutTime(ctx, limit)
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

func (p *ParserController) getGamesWithoutTime(ctx context.Context, limit int64) ([]games.GameRecord, error) {
	var records []*core.Record
	err := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		WithContext(ctx).
		Where(dbx.HashExp{"campaign_time": -1}).
		Limit(limit).
		All(&records)
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
