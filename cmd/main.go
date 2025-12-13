package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/games/parser"
	stracker "adventuria/internal/adventuria/stream-tracker"
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
		gamesParser, err := parser.NewGamesParser()
		if err == nil {
			adventuria.PocketBase.Cron().MustAdd("games_parser", "0 0 1 * *", func() {
				gamesParser.Parse(game.Context())
			})
		}

		st, err := stracker.NewStreamTracker()
		if err != nil {
			adventuria.PocketBase.Logger().Error("Failed to initialize stream tracker", "error", err)
		} else {
			if err = st.Start(game.Context()); err != nil {
				adventuria.PocketBase.Logger().Error("Failed to start stream tracker", "error", err)
			}
		}

		http.Route(game, se.Router)

		return se.Next()
	}); err != nil {
		log.Fatal(err)
	}
}
