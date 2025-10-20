package http

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/etag"
	"os"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func Route(game adventuria.Game, router *router.Router[*core.RequestEvent]) {
	handlers := New(game)

	gs := router.Group("")
	gs.BindFunc(apis.WrapStdMiddleware(etag.Etag))
	gs.GET("/{path...}", apis.Static(os.DirFS("./static"), false))

	g := router.Group("/api")

	g.GET("/timer/left/{userId}", handlers.GetTimeLeftByUserHandler)

	ga := g.Group("")
	ga.Bind(apis.RequireAuth())

	ga.GET("/timer/left", handlers.GetTimeLeftHandler)

	gab := ga.Group("")
	gab.Bind(adventuria.GameSettings.CheckActionsBlock())

	gab.POST("/roll", handlers.RollHandler)

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
}
