package igdb

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ParserController struct {
	parser games.Parser
}

func NewParserController() *ParserController {
	p := &ParserController{
		parser: NewParser(
			adventuria.GameSettings.TwitchClientID(),
			adventuria.GameSettings.TwitchClientSecret(),
			adventuria.GameSettings.IGDBParseSettings(),
		),
	}

	return p

	//adventuria.PocketBase.Cron().MustAdd("games_parser", "0 3 * * 0", p.ParseGames)
}

func (p *ParserController) Parse() {
	p.parseCompanies()
	p.parsePlatforms()
	p.parseGames()
}

func (p *ParserController) parseGames() {
	ch, err := p.parser.ParseGames()
	if err != nil {
		// TODO error handling
	}

	for games := range ch {
		for _, game := range games {
			// TODO save to PocketBase
			_ = game
		}
	}
}

func (p *ParserController) parsePlatforms() {
	ch, err := p.parser.ParsePlatforms()
	if err != nil {
		// TODO error handling
	}

	for platforms := range ch {
		for _, platform := range platforms {
			// TODO save to PocketBase
			_ = platform
		}
	}
}

func (p *ParserController) parseCompanies() {
	ch, err := p.parser.ParseCompanies()
	if err != nil {
		// TODO error handling
	}

	for companies := range ch {
		for _, company := range companies {
			// TODO save to PocketBase
			_ = company
		}
	}
}

func (p *ParserController) platformsToIds(platformIdsDb []uint64) ([]string, error) {
	col, err := adventuria.GameCollections.Get(adventuria.CollectionPlatforms)
	if err != nil {
		return nil, err
	}

	records := make([]*core.Record, len(platformIdsDb))
	err = adventuria.PocketBase.
		RecordQuery(col).
		Where(
			dbx.Or(
				dbx.HashExp{"igdb_id": platformIdsDb},
			),
		).
		All(&records)
	if err != nil {
		return nil, err
	}

	platformIds := make([]string, len(records))
	for i, record := range records {
		platformIds[i] = record.Id
	}

	return platformIds, nil
}
