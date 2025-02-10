package main

import (
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

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./static"), false))

		g := se.Router.Group("/api")
		g.Bind(apis.RequireAuth())

		g.POST("/roll", handlers.RollHandler)
		g.POST("/choose-game", handlers.ChooseGameHandler)

		g.GET("/get-last-action", handlers.GetLastActionHandler)
		g.GET("/game-result", handlers.GameResultHandler)

		g.POST("/reroll", handlers.RerollHandler)
		g.POST("/drop", handlers.DropHandler)
		g.POST("/done", handlers.DoneHandler)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
