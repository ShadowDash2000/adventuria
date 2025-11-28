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

func NewGamesParser(ctx context.Context) (*GamesParser, error) {
	igdbParser, err := igdb.New(ctx)
	if err != nil {
		log.Printf("Failed to initialize igdb parser: %v", err)
		return nil, err
	}

	steamParser, err := steam.New(ctx)
	if err != nil {
		log.Printf("Failed to initialize steam parser: %v", err)
		return nil, err
	}

	hltbParser, err := hltb.New(ctx)
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

func (p *GamesParser) Parse() {
	adventuria.PocketBase.Logger().Info("IGDB parser started")
	p.igdbParser.Parse(500)
	adventuria.PocketBase.Logger().Info("IGDB parser finished")

	adventuria.PocketBase.Logger().Info("Steam parser started")
	p.steamParser.Parse(100)
	adventuria.PocketBase.Logger().Info("Steam parser finished")

	adventuria.PocketBase.Logger().Info("HLTB parser started")
	p.hltbParser.Parse(100)
	adventuria.PocketBase.Logger().Info("HLTB parser finished")
}
