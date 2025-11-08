package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/games/hltb"
	"adventuria/internal/adventuria/games/igdb"
	"adventuria/internal/adventuria/games/steam"
	"adventuria/internal/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	game := adventuria.New()

	actions.WithBaseActions()
	effects.WithBaseEffects()
	cells.WithBaseCells()

	game.OnServe(func(se *core.ServeEvent) error {
		igdbParser, err := igdb.New()
		if err != nil {
			log.Printf("Failed to initialize igdb parser: %v", err)
		} else {
			adventuria.PocketBase.Cron().MustAdd("igdb_parser", "0 0 1 * *", igdbParser.Parse)
		}

		steamParser, err := steam.New()
		if err != nil {
			log.Printf("Failed to initialize steam parser: %v", err)
		} else {
			adventuria.PocketBase.Cron().MustAdd("steam_prices_parser", "0 0 1 * *", steamParser.Parse)
		}

		hltbParser, err := hltb.New()
		if err != nil {
			log.Printf("Failed to initialize hltb parser: %v", err)
		} else {
			adventuria.PocketBase.Cron().MustAdd("hltb_parser", "0 0 1 * *", hltbParser.Parse)
		}

		http.Route(game, se.Router)

		return se.Next()
	})

	if err := game.Start(); err != nil {
		log.Fatal(err)
	}
}
