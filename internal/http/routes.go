package http

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func Route(game game, router *router.Router[*core.RequestEvent]) {
	handlers := New(game)

	g := router.Group("/api")

	g.GET("/event-stats", handlers.EventStats)
	g.GET("/current-season", handlers.CurrentSeason)
	g.GET("/event-ended", handlers.IsEventEnded)

	ga := g.Group("")
	ga.Bind(apis.RequireAuth())

	gab := ga.Group("")
	gab.BindFunc(func(e *core.RequestEvent) error {
		isBlocked, err := handlers.Game.IsActionsBlocked(e.Request.Context())
		if err != nil {
			return RespondWithError(e, err)
		}

		if isBlocked {
			return e.JSON(http.StatusForbidden, result{
				Success: false,
				Data:    nil,
				Error:   "all_actions_blocked",
				Message: "All actions are blocked",
			})
		}

		return e.Next()
	})

	gab.POST("/roll", handlers.RollHandler)

	gab.POST("/update-action", handlers.UpdateReviewHandler)
	gab.GET("/available-actions", handlers.GetAvailableActions)
	gab.GET("/action-view", handlers.GetActionView)

	gab.POST("/reroll", handlers.RerollHandler)
	gab.POST("/drop", handlers.DropHandler)
	gab.POST("/done", handlers.DoneHandler)

	gab.POST("/generate-wheel", handlers.GenerateWheelHandler)
	gab.POST("/roll-wheel", handlers.RollWheelHandler)
	gab.POST("/roll-item", handlers.RollItemHandler)
	gab.POST("/roll-item-on-cell", handlers.RollItemOnCellHandler)

	gab.POST("/buy-item", handlers.BuyItemHandler)
	gab.POST("/refresh-shop", handlers.RefreshShopHandler)

	gab.POST("/use-item", handlers.UseItemHandler)
	gab.POST("/drop-item", handlers.DropItemHandler)
	gab.POST("/effect-view", handlers.GetEffectView)
}
