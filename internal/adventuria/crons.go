package adventuria

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
)

func (g *Game) registerCrons(ctx context.Context, pb core.App, registry *Registry) {
	pb.Cron().MustAdd("games_parser", "0 0 1 * *", func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		settings, err := g.settings.GetFirst(ctx)
		if err != nil {
			return
		}

		unsub := g.onKillParser().BindFunc(func(ctx context.Context, e *onKillParserEvent) error {
			cancel()
			return e.Next()
		})
		defer unsub()

		if !settings.DisableSteamParser() {
			pb.Logger().Info("Started SteamSpy parser")

			err = registry.SteamSpy().Parse(ctx)
			if err != nil {
				pb.Logger().Error("SteamSpy parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished SteamSpy parser")
		}

		if !settings.DisableCheapsharkParser() {
			pb.Logger().Info("Started CheapShark parser")

			err = registry.CheapShark().Parse(ctx)
			if err != nil {
				pb.Logger().Error("CheapShark parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished CheapShark parser")
		}

		if !settings.DisableHltbParser() {
			pb.Logger().Info("Started HLTB parser")

			err = registry.HLTB().Parse(ctx)
			if err != nil {
				pb.Logger().Error("HLTB parser failed", "error", err)
				return
			}

			pb.Logger().Info("Finished HLTB parser")
		}

		if !settings.DisableIgdbParser() {
			err = registry.IGDB().ParsePlatforms(ctx, 500)
			if err != nil {
				return
			}
			err = registry.IGDB().ParseGenres(ctx, 500)
			if err != nil {
				return
			}
			err = registry.IGDB().ParseGameTypes(ctx, 500)
			if err != nil {
				return
			}
			if !settings.DisableIgdbGamesParser() {
				err = registry.IGDB().ParseGames(ctx, settings.IgdbFilter().Build(), 500)
				if err != nil {
					return
				}
			}
		}
	})

	pb.Cron().MustAdd("cell_events_scheduler", "*/1 * * * *", func() {
		err := g.cellEvents.CheckEventsSchedules(ctx)
		if err != nil {
			pb.Logger().Error("Failed to run cell events scheduler", "error", err)
			return
		}
	})

	pb.Cron().MustAdd("update_cells_average_campaign_time", "0 4 * * *", func() {
		ids, err := g.cells.GetAllIDs(ctx)
		if err != nil {
			return
		}

		for _, id := range ids {
			campaignTime, err := g.cells.GetAverageCampaignTimeByID(ctx, id)
			if err != nil {
				pb.Logger().Error("Failed to get average campaign time", "error", err)
				return
			}

			err = g.cells.UpdateAverageCampaignTimeByID(ctx, id, campaignTime)
			if err != nil {
				pb.Logger().Error("Failed to update average campaign time", "error", err)
				return
			}
		}
	})
}
