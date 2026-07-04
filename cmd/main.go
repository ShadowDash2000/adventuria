package main

import (
	"adventuria/internal/adventuria_new"
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

	_, err := adventuria_new.Start(func(game *adventuria_new.Game, se *core.ServeEvent) error {
		/*gamesParser, err := parser.NewGamesParser()
		if err == nil {
			se.App.Cron().MustAdd("games_parser", "0 0 1 * *", func() {
				gamesParser.Parse(context.Background())
			})
			se.App.Cron().MustAdd("refresh_hltb_time", "0 0 1 * *", func() {
				gamesParser.RefreshHltbTime(context.Background())
			})
		}*/

		http.Route(game, se.Router)

		return se.Next()
	})
	if err != nil {
		log.Fatal(err)
	}
}
