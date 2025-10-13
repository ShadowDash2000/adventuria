package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/http"
	"log"

	"github.com/pocketbase/pocketbase/core"
)

func main() {
	game := adventuria.New()

	actions.WithBaseActions()
	effects.WithBaseEffects()
	cells.WithBaseCells()

	game.OnServe(func(se *core.ServeEvent) error {
		http.Route(game, se.Router)

		return se.Next()
	})

	if err := game.Start(); err != nil {
		log.Fatal(err)
	}
}
