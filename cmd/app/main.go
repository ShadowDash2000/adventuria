package main

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/dlc/pack1"
	"adventuria/internal/adventuria/dlc/pack2"
	"adventuria/internal/http/handlers/v1"
	"adventuria/pkg/etag"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"os"

	_ "adventuria/migrations"
)

func main() {
	app := pocketbase.New()

	game := adventuria.New(app)

	handlers := handlers.New(game)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		game.Init()

		game = pack1.WithItemPack1(game)
		game = pack2.WithItemPack2(game)

		// TODO include parser only after twitch auth is set
		/*games.NewParser(
			adventuria.GameSettings.TwitchClientID(),
			adventuria.GameSettings.TwitchClientSecret(),
			adventuria.GameSettings.IGDBParseSettings(),
			adventuria.GameCollections,
			adventuria.GameApp,
		)*/

		gs := se.Router.Group("")
		gs.BindFunc(apis.WrapStdMiddleware(etag.Etag))
		gs.GET("/{path...}", apis.Static(os.DirFS("./static"), false))

		g := se.Router.Group("/api")

		g.GET("/timer/left/{userId}", handlers.GetTimeLeftByUserHandler)

		ga := g.Group("")
		ga.Bind(apis.RequireAuth())

		ga.GET("/timer/left", handlers.GetTimeLeftHandler)

		gab := ga.Group("")
		gab.Bind(adventuria.GameSettings.CheckActionsBlock())

		gab.POST("/roll", handlers.RollHandler)

		gab.GET("/get-next-step-type", handlers.GetNextStepTypeHandler)
		gab.GET("/get-last-action", handlers.GetLastActionHandler)
		gab.GET("/get-roll-effects", handlers.GetRollEffectsHandler)

		gab.POST("/update-action", handlers.UpdateActionHandler)

		gab.POST("/reroll", handlers.RerollHandler)
		gab.POST("/drop", handlers.DropHandler)
		gab.POST("/done", handlers.DoneHandler)

		gab.POST("/roll-wheel", handlers.RollWheelHandler)
		gab.POST("/roll-item", handlers.RollItemHandler)

		gab.POST("/use-item", handlers.UseItemHandler)
		gab.POST("/drop-item", handlers.DropItemHandler)

		gab.POST("/timer/start", handlers.StartTimerHandler)
		gab.POST("/timer/stop", handlers.StopTimerHandler)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
