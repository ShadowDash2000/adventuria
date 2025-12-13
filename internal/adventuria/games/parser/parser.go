package parser

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games/hltb"
	"adventuria/internal/adventuria/games/igdb"
	"adventuria/internal/adventuria/games/steam"
	"context"
	"log"
)

type GamesParser struct {
	igdbParser  *igdb.ParserController
	steamParser *steam.ParserController
	hltbParser  *hltb.ParserController
}

func NewGamesParser() (*GamesParser, error) {
	igdbParser, err := igdb.New()
	if err != nil {
		log.Printf("Failed to initialize igdb parser: %v", err)
		return nil, err
	}

	steamParser, err := steam.New()
	if err != nil {
		log.Printf("Failed to initialize steam parser: %v", err)
		return nil, err
	}

	hltbParser, err := hltb.New()
	if err != nil {
		log.Printf("Failed to initialize hltb parser: %v", err)
		return nil, err
	}

	return &GamesParser{
		igdbParser:  igdbParser,
		steamParser: steamParser,
		hltbParser:  hltbParser,
	}, nil
}

func (p *GamesParser) Parse(ctx context.Context) {
	if !adventuria.GameSettings.DisableIGDBParser() {
		adventuria.PocketBase.Logger().Info("IGDB parser started")
		p.igdbParser.Parse(ctx, 500)
		adventuria.PocketBase.Logger().Info("IGDB parser finished")
	}

	if !adventuria.GameSettings.DisableSteamParser() {
		adventuria.PocketBase.Logger().Info("Steam parser started")
		p.steamParser.Parse(ctx)
		adventuria.PocketBase.Logger().Info("Steam parser finished")
	}

	if !adventuria.GameSettings.DisableHLTBParser() {
		adventuria.PocketBase.Logger().Info("HLTB parser started")
		p.hltbParser.Parse(ctx, 100)
		adventuria.PocketBase.Logger().Info("HLTB parser finished")
	}
}
