package hltb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"context"
	"math"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser *Parser
}

func New() *ParserController {
	return &ParserController{
		parser: NewParser(),
	}
}

func (p *ParserController) Parse(ctx context.Context) {
	if err := p.parseGames(ctx); err != nil {
		adventuria.PocketBase.Logger().Error("Failed to parse games", "error", err)
	}
}

func (p *ParserController) parseGames(ctx context.Context) error {
	games, err := p.parser.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	adventuria.PocketBase.Logger().Info("hltb.parseGames", "games_count", len(games))

	gamesProcessed := 0
	gamesSkipped := 0
	for _, game := range games {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if ok := p.isGameExist(ctx, game.GameID); ok {
			gamesSkipped++
			continue
		}

		gameRecord := NewHowLongToBeatRecordFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(schema.CollectionHowLongToBeat),
			),
		)
		gameRecord.SetIdDb(game.GameID)
		gameRecord.SetName(game.GameName)
		gameRecord.SetYear(game.ReleaseWorld)
		gameRecord.SetCampaign(math.Round(float64(game.CompMain) / 3600))
		if err = adventuria.PocketBase.Save(gameRecord.ProxyRecord()); err != nil {
			return err
		}

		gamesProcessed++
	}

	adventuria.PocketBase.Logger().Info("hltb.parseGames",
		"games_processed", gamesProcessed,
		"games_skipped", gamesSkipped,
	)

	return nil
}

func (p *ParserController) isGameExist(ctx context.Context, id int) bool {
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(schema.CollectionHowLongToBeat)).
		WithContext(ctx).
		Where(dbx.HashExp{"id_db": id}).
		One(&core.Record{})
	return err == nil
}
