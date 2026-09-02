package http

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/http/completed_activities"
	"adventuria/internal/http/event_stats"
	"adventuria/internal/http/response"
	"adventuria/internal/http/update_review"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func Route(game *adventuria.Game, registry *adventuria.Registry, router *router.Router[*core.RequestEvent]) {
	handlers := New(game, registry.Settings())
	completedActivities := completed_activities.NewHandler(registry.CompletedActivities())
	eventStats := event_stats.NewHandler(registry.EventStats())
	updateReview := update_review.NewHandler(registry.Reviews())

	g := router.Group("/api")

	g.GET("/event-stats", eventStats.EventStats)
	g.GET("/current-season", handlers.CurrentSeason)
	g.GET("/event-ended", handlers.IsEventEnded)
	g.GET("/completed-activities", completedActivities.GetCompletedActivitiesByCellID)

	ga := g.Group("")
	ga.Bind(apis.RequireAuth())

	gab := ga.Group("")
	gab.BindFunc(func(e *core.RequestEvent) error {
		err := game.IsActionsBlocked(e.Request.Context())
		if err != nil {
			return response.Error(e, err)
		}

		isDisabled, err := registry.Players().IsDisabled(e.Request.Context(), e.Auth.Id)
		if err != nil {
			return response.Error(e, err)
		}
		if isDisabled {
			return response.Error(e, errs.ErrPlayerIsDisabled)
		}

		return e.Next()
	})

	gab.POST("/start", handlers.StartHandler)

	gab.POST("/roll", handlers.RollHandler)

	gab.POST("/update-review", updateReview.UpdateReviewByActionID)
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
	gab.POST("/coins-for-item", handlers.CoinsForItemHandler)

	gab.POST("/use-item", handlers.UseItemHandler)
	gab.POST("/drop-item", handlers.DropItemHandler)
	gab.POST("/effect-view", handlers.GetEffectView)
}
