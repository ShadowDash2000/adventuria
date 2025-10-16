package steam

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
	"context"
	"errors"
	"os"
)

type ParserController struct {
	parser *Parser
}

func New() (*ParserController, error) {
	steamApiKey, ok := os.LookupEnv("STEAM_API_KEY")
	if !ok {
		return nil, errors.New("steam: STEAM_API_KEY not found")
	}

	p := &ParserController{
		parser: NewParser(steamApiKey),
	}

	return p, nil
}

func (p *ParserController) Parse() {
	const limit = 100
	for {
		games, err := p.getSteamAppsWithoutPrice(limit)
		if err != nil {
			adventuria.PocketBase.Logger().Error("Failed to get games", "error", err)
			return
		}

		if len(games) == 0 {
			break
		}

		err = p.parser.ParsePrices(context.Background(), games)
		if err != nil {
			adventuria.PocketBase.Logger().Error("Failed to parse prices", "error", err)
			return
		}

		for _, game := range games {
			err = adventuria.PocketBase.Save(game.ProxyRecord())
			if err != nil {
				adventuria.PocketBase.Logger().Error("Failed to save game", "error", err)
				return
			}
		}
	}
}

func (p *ParserController) getSteamAppsWithoutPrice(limit int) ([]games.GameRecord, error) {
	records, err := adventuria.PocketBase.FindRecordsByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		"platforms.id_db ?= 6 && steam_app_id != 0 && steam_app_price = -1",
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
