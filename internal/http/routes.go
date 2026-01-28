package http

import (
	"adventuria/internal/adventuria"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func Route(game *adventuria.Game, router *router.Router[*core.RequestEvent]) {
	handlers := New(game)

	g := router.Group("/api")

	timer := g.Group("/timer")
	timer.GET("/left/{userId}", handlers.GetTimeLeftByUserHandler)
	timer.GET("/left", handlers.GetTimeLeftHandler)

	timerA := timer.Group("")
	timerA.Bind(apis.RequireAuth())
	timerA.Bind(adventuria.GameSettings.CheckActionsBlock())
	timerA.POST("/start", handlers.StartTimerHandler)
	timerA.POST("/stop", handlers.StopTimerHandler)

	ga := g.Group("")
	ga.Bind(apis.RequireAuth())

	gab := ga.Group("")
	gab.Bind(adventuria.GameSettings.CheckActionsBlock())

	gab.POST("/roll", handlers.RollHandler)

	gab.POST("/update-action", handlers.UpdateActionHandler)
	gab.GET("/available-actions", handlers.GetAvailableActions)

	gab.POST("/reroll", handlers.RerollHandler)
	gab.POST("/drop", handlers.DropHandler)
	gab.POST("/done", handlers.DoneHandler)

	gab.POST("/roll-wheel", handlers.RollWheelHandler)
	gab.POST("/roll-item", handlers.RollItemHandler)
	gab.POST("/roll-item-on-cell", handlers.RollItemOnCellHandler)
	gab.POST("/buy-item", handlers.BuyItemHandler)

	gab.POST("/use-item", handlers.UseItemHandler)
	gab.POST("/drop-item", handlers.DropItemHandler)
	gab.POST("/get-item-effect-variants", handlers.GetItemEffectVariants)
}
