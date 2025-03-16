package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/http/handlers/v1"
	_ "adventuria/migrations"
	"adventuria/pkg/etag"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"os"
)

func main() {
	app := pocketbase.New()

	game := adventuria.NewGame(app.App)

	handlers := handlers.New(game)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		game.Init()

		gs := se.Router.Group("")
		gs.BindFunc(apis.WrapStdMiddleware(etag.Etag))
		gs.GET("/{path...}", apis.Static(os.DirFS("./static"), false))

		g := se.Router.Group("/api")
		g.Bind(game.GC.Settings.CheckActionsBlock(), apis.RequireAuth())

		g.POST("/roll", handlers.RollHandler)
		g.POST("/choose-game", handlers.ChooseGameHandler)

		g.GET("/get-next-step-type", handlers.GetNextStepTypeHandler)
		g.GET("/get-last-action", handlers.GetLastActionHandler)
		g.GET("/get-roll-effects", handlers.GetRollEffectsHandler)

		g.POST("/reroll", handlers.RerollHandler)
		g.POST("/drop", handlers.DropHandler)
		g.POST("/done", handlers.DoneHandler)

		g.POST("/roll-cell", handlers.RollCellHandler)
		g.POST("/roll-wheel-preset", handlers.RollWheelPresetHandler)
		g.POST("/roll-item", handlers.RollItemHandler)

		g.POST("/use-item", handlers.UseItemHandler)
		g.POST("/drop-item", handlers.DropItemHandler)

		g.POST("/timer/start", handlers.StartTimerHandler)
		g.POST("/timer/stop", handlers.StopTimerHandler)
		g.GET("/timer/left", handlers.GetTimeLeftHandler)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
