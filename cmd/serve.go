package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/games/igdb"
	"adventuria/internal/adventuria/games/steam"
	"adventuria/internal/http"
	"adventuria/pkg/config"
	"log"

	"github.com/pocketbase/pocketbase/core"
)

func main() {
	config.LoadEnv()

	game := adventuria.New()

	actions.WithBaseActions()
	effects.WithBaseEffects()
	cells.WithBaseCells()

	game.OnServe(func(se *core.ServeEvent) error {
		igdbParser, err := igdb.New()
		if err != nil {
			return err
		}
		adventuria.PocketBase.Cron().MustAdd("igdb_parser", "0 0 1 * *", igdbParser.Parse)

		steamParser, err := steam.New()
		if err != nil {
			return err
		}
		adventuria.PocketBase.Cron().MustAdd("steam_prices_parser", "0 0 1 * *", steamParser.Parse)

		http.Route(game, se.Router)

		return se.Next()
	})

	if err := game.Start(); err != nil {
		log.Fatal(err)
	}
}
