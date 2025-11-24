package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/games/parser"
	"adventuria/internal/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase/core"

	_ "adventuria/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	game := adventuria.New()

	actions.WithBaseActions()
	effects.WithBaseEffects()
	cells.WithBaseCells()

	if err := game.Start(func(se *core.ServeEvent) error {
		gamesParser, err := parser.NewGamesParser(game.Context())
		if err == nil {
			adventuria.PocketBase.Cron().MustAdd("games_parser", "0 0 1 * *", gamesParser.Parse)
		}

		http.Route(game, se.Router)

		return se.Next()
	}); err != nil {
		log.Fatal(err)
	}
}
