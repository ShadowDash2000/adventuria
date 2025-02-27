package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/http/handlers/v1"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"os"
)

func main() {
	//cfg := config.MustLoad()

	app := pocketbase.New()

	game := adventuria.NewGame(app.App)

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return game.Init()
	})

	handlers := handlers.New(game)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./static"), false))

		g := se.Router.Group("/api")
		g.Bind(apis.RequireAuth())

		g.POST("/roll", handlers.RollHandler)
		g.POST("/choose-game", handlers.ChooseGameHandler)

		g.GET("/get-next-step-type", handlers.GetNextStepTypeHandler)
		g.GET("/game-result", handlers.GameResultHandler)
		g.GET("/get-roll-effects", handlers.GetRollEffectsHandler)

		g.POST("/reroll", handlers.RerollHandler)
		g.POST("/drop", handlers.DropHandler)
		g.POST("/done", handlers.DoneHandler)

		g.POST("/roll-cell", handlers.RollCellHandler)
		g.POST("/roll-movie", handlers.RollMovieHandler)
		g.POST("/roll-item", handlers.RollItemHandler)
		g.POST("/roll-big-win", handlers.RollBigWinHandler)
		g.POST("/roll-developer", handlers.RollDeveloperHandler)

		g.POST("/use-item", handlers.UseItemHandler)
		g.POST("/drop-item", handlers.DropItemHandler)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
