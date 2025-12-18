package parser

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games/cheapshark"
	"adventuria/internal/adventuria/games/hltb"
	"adventuria/internal/adventuria/games/igdb"
	"adventuria/internal/adventuria/games/steam"
	"context"
	"log"
)

type GamesParser struct {
	igdbParser       *igdb.ParserController
	steamParser      *steam.ParserController
	hltbParser       *hltb.ParserController
	cheapsharkParser *cheapshark.ParserController
}

func NewGamesParser() (*GamesParser, error) {
	igdbParser, err := igdb.New()
	if err != nil {
		log.Printf("Failed to initialize igdb parser: %v", err)
		return nil, err
	}

	return &GamesParser{
		igdbParser:       igdbParser,
		steamParser:      steam.New(),
		hltbParser:       hltb.New(),
		cheapsharkParser: cheapshark.New(),
	}, nil
}

func (p *GamesParser) Parse(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	unsub := adventuria.GameSettings.OnKillParser().BindFunc(func(e *adventuria.OnKillParserEvent) error {
		cancel()
		return e.Next()
	})
	defer unsub()

	if !adventuria.GameSettings.DisableSteamParser() {
		adventuria.PocketBase.Logger().Info("Steam parser started")
		p.steamParser.Parse(ctx)
		adventuria.PocketBase.Logger().Info("Steam parser finished")
	}

	if !adventuria.GameSettings.DisableCheapsharkParser() {
		adventuria.PocketBase.Logger().Info("Cheapshark parser started")
		p.cheapsharkParser.Parse(ctx)
		adventuria.PocketBase.Logger().Info("Cheapshark parser finished")
	}

	if !adventuria.GameSettings.DisableHLTBParser() {
		adventuria.PocketBase.Logger().Info("HLTB scraper parser started")
		p.hltbParser.Parse(ctx)
		adventuria.PocketBase.Logger().Info("HLTB scraper parser finished")
	}

	if !adventuria.GameSettings.DisableIGDBParser() {
		adventuria.PocketBase.Logger().Info("IGDB parser started")
		p.igdbParser.Parse(ctx, 500)
		adventuria.PocketBase.Logger().Info("IGDB parser finished")
	}
}
