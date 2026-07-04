package adventuria_new

import (
	"adventuria/internal/adventuria_new/actions"
	rollWheelRepo "adventuria/internal/adventuria_new/actions/custom/roll_wheel/repository"
	actionsRepo "adventuria/internal/adventuria_new/actions/repository"
	"adventuria/internal/adventuria_new/activities"
	activitiesRepo "adventuria/internal/adventuria_new/activities/repository"
	"adventuria/internal/adventuria_new/activity_filters"
	activityFiltersRepo "adventuria/internal/adventuria_new/activity_filters/repository"
	"adventuria/internal/adventuria_new/board"
	"adventuria/internal/adventuria_new/cells"
	cellsRepo "adventuria/internal/adventuria_new/cells/repository"
	"adventuria/internal/adventuria_new/effects"
	effectsRepo "adventuria/internal/adventuria_new/effects/repository"
	"adventuria/internal/adventuria_new/event_stats"
	eventStatsRepo "adventuria/internal/adventuria_new/event_stats/repository"
	"adventuria/internal/adventuria_new/genres"
	genresRepo "adventuria/internal/adventuria_new/genres/repository"
	"adventuria/internal/adventuria_new/inventories"
	inventoriesRepo "adventuria/internal/adventuria_new/inventories/repository"
	"adventuria/internal/adventuria_new/items"
	itemsRepo "adventuria/internal/adventuria_new/items/repository"
	"adventuria/internal/adventuria_new/outboxes"
	outboxesRepo "adventuria/internal/adventuria_new/outboxes/repository"
	"adventuria/internal/adventuria_new/player_progress"
	progressRepo "adventuria/internal/adventuria_new/player_progress/repository"
	"adventuria/internal/adventuria_new/players"
	playersRepo "adventuria/internal/adventuria_new/players/repository"
	"adventuria/internal/adventuria_new/reviews"
	reviewsRepo "adventuria/internal/adventuria_new/reviews/repository"
	"adventuria/internal/adventuria_new/seasons"
	seasonsRepo "adventuria/internal/adventuria_new/seasons/repository"
	"adventuria/internal/adventuria_new/settings"
	settingsRepo "adventuria/internal/adventuria_new/settings/repository"
	"adventuria/internal/adventuria_new/stream_tracker"
	streamTrackerRepo "adventuria/internal/adventuria_new/stream_tracker/repository"
	"adventuria/internal/adventuria_new/worlds"
	worldsRepo "adventuria/internal/adventuria_new/worlds/repository"

	"github.com/pocketbase/pocketbase/core"
)

type Registry struct {
	pb core.App

	// repos
	seasonsRepo          *seasonsRepo.Repository
	settingsRepo         *settingsRepo.Repository
	cellsRepo            *cellsRepo.Repository
	worldsRepo           *worldsRepo.Repository
	actionsRepo          *actionsRepo.Repository
	progressRepo         *progressRepo.Repository
	playersRepo          *playersRepo.Repository
	inventoriesRepo      *inventoriesRepo.Repository
	effectsRepo          *effectsRepo.Repository
	activitiesRepo       *activitiesRepo.Repository
	activityFiltersRepo  *activityFiltersRepo.Repository
	itemsRepo            *itemsRepo.Repository
	genresRepo           *genresRepo.Repository
	reviewsRepo          *reviewsRepo.Repository
	rollWheelRepo        *rollWheelRepo.Repository
	outboxesRepo         *outboxesRepo.Repository
	relationRepo         *activitiesRepo.RelationRepository
	eventStatsRepo       *eventStatsRepo.Repository
	eventStatsCachedRepo *eventStatsRepo.CachedRepository
	streamTrackerRepo    *streamTrackerRepo.Repository

	// services
	seasons         *seasons.Seasons
	settings        *settings.Settings
	worlds          *worlds.Worlds
	cells           *cells.Cells
	actions         *actions.Actions
	progress        *player_progress.PlayerProgress
	players         *players.Players
	effects         *effects.Effects
	inventories     *inventories.Inventories
	activities      *activities.Activities
	activityFilters *activity_filters.ActivityFilters
	items           *items.Items
	board           *board.Board
	genres          *genres.Genres
	reviews         *reviews.Reviews
	outboxes        *outboxes.Outboxes
	eventStats      *event_stats.EventStats
	streamTracker   *stream_tracker.StreamTracker
}

func NewRegistry(pb core.App) *Registry {
	return &Registry{pb: pb}
}

func (r *Registry) SeasonsRepo() *seasonsRepo.Repository {
	if r.seasonsRepo == nil {
		r.seasonsRepo = seasonsRepo.NewRepository(r.pb)
	}
	return r.seasonsRepo
}

func (r *Registry) SettingsRepo() *settingsRepo.Repository {
	if r.settingsRepo == nil {
		r.settingsRepo = settingsRepo.NewRepository(r.pb)
	}
	return r.settingsRepo
}

func (r *Registry) CellsRepo() *cellsRepo.Repository {
	if r.cellsRepo == nil {
		r.cellsRepo = cellsRepo.NewRepository(r.pb)
	}
	return r.cellsRepo
}

func (r *Registry) WorldsRepo() *worldsRepo.Repository {
	if r.worldsRepo == nil {
		r.worldsRepo = worldsRepo.NewRepository(r.pb)
	}
	return r.worldsRepo
}

func (r *Registry) ActionsRepo() *actionsRepo.Repository {
	if r.actionsRepo == nil {
		r.actionsRepo = actionsRepo.NewRepository(r.pb)
	}
	return r.actionsRepo
}

func (r *Registry) PlayerProgressRepo() *progressRepo.Repository {
	if r.progressRepo == nil {
		r.progressRepo = progressRepo.NewRepository(r.pb)
	}
	return r.progressRepo
}

func (r *Registry) PlayersRepo() *playersRepo.Repository {
	if r.playersRepo == nil {
		r.playersRepo = playersRepo.NewRepository(r.pb)
	}
	return r.playersRepo
}

func (r *Registry) InventoriesRepo() *inventoriesRepo.Repository {
	if r.inventoriesRepo == nil {
		r.inventoriesRepo = inventoriesRepo.NewRepository(r.pb)
	}
	return r.inventoriesRepo
}

func (r *Registry) EffectsRepo() *effectsRepo.Repository {
	if r.effectsRepo == nil {
		r.effectsRepo = effectsRepo.NewRepository(r.pb)
	}
	return r.effectsRepo
}

func (r *Registry) ActivitiesRepo() *activitiesRepo.Repository {
	if r.activitiesRepo == nil {
		r.activitiesRepo = activitiesRepo.NewRepository(r.pb)
	}
	return r.activitiesRepo
}

func (r *Registry) ActivityFiltersRepo() *activityFiltersRepo.Repository {
	if r.activityFiltersRepo == nil {
		r.activityFiltersRepo = activityFiltersRepo.NewRepository(r.pb)
	}
	return r.activityFiltersRepo
}

func (r *Registry) ItemsRepo() *itemsRepo.Repository {
	if r.itemsRepo == nil {
		r.itemsRepo = itemsRepo.NewRepository(r.pb)
	}
	return r.itemsRepo
}

func (r *Registry) GenresRepo() *genresRepo.Repository {
	if r.genresRepo == nil {
		r.genresRepo = genresRepo.NewRepository(r.pb)
	}
	return r.genresRepo
}

func (r *Registry) ReviewsRepo() *reviewsRepo.Repository {
	if r.reviewsRepo == nil {
		r.reviewsRepo = reviewsRepo.NewRepository(r.pb)
	}
	return r.reviewsRepo
}

func (r *Registry) RollWheelRepo() *rollWheelRepo.Repository {
	if r.rollWheelRepo == nil {
		r.rollWheelRepo = rollWheelRepo.NewRepository(r.pb)
	}
	return r.rollWheelRepo
}

func (r *Registry) OutboxesRepo() *outboxesRepo.Repository {
	if r.outboxesRepo == nil {
		r.outboxesRepo = outboxesRepo.NewRepository(r.pb)
	}
	return r.outboxesRepo
}

func (r *Registry) RelationRepo() *activitiesRepo.RelationRepository {
	if r.relationRepo == nil {
		r.relationRepo = activitiesRepo.NewRelationRepository(r.pb)
	}
	return r.relationRepo
}

func (r *Registry) EventStatsRepo() *eventStatsRepo.Repository {
	if r.eventStatsRepo == nil {
		r.eventStatsRepo = eventStatsRepo.NewRepository(r.pb)
	}
	return r.eventStatsRepo
}

func (r *Registry) EventStatsCachedRepo() *eventStatsRepo.CachedRepository {
	if r.eventStatsCachedRepo == nil {
		r.eventStatsCachedRepo = eventStatsRepo.NewCachedRepository(r.EventStatsRepo())
	}
	return r.eventStatsCachedRepo
}

func (r *Registry) StreamTrackerRepo() *streamTrackerRepo.Repository {
	if r.streamTrackerRepo == nil {
		r.streamTrackerRepo = streamTrackerRepo.NewRepository(r.pb)
	}
	return r.streamTrackerRepo
}

func (r *Registry) Seasons() *seasons.Seasons {
	if r.seasons == nil {
		r.seasons = seasons.NewSeasons(r.SeasonsRepo())
	}
	return r.seasons
}

func (r *Registry) Settings() *settings.Settings {
	if r.settings == nil {
		r.settings = settings.NewSettings(r.SettingsRepo(), r.Seasons())
	}
	return r.settings
}

func (r *Registry) Worlds() *worlds.Worlds {
	if r.worlds == nil {
		r.worlds = worlds.NewWorlds(r.WorldsRepo(), r.Effects())
	}
	return r.worlds
}

func (r *Registry) Cells() *cells.Cells {
	if r.cells == nil {
		r.cells = cells.NewCells(r.CellsRepo())
	}
	return r.cells
}

func (r *Registry) Actions() *actions.Actions {
	if r.actions == nil {
		r.actions = actions.NewActions(r.ActionsRepo(), r.Worlds(), r.Cells())
	}
	return r.actions
}

func (r *Registry) PlayerProgress() *player_progress.PlayerProgress {
	if r.progress == nil {
		r.progress = player_progress.NewPlayerProgress(r.PlayerProgressRepo(), r.Worlds())
	}
	return r.progress
}

func (r *Registry) Players() *players.Players {
	if r.players == nil {
		r.players = players.NewPlayers(
			r.PlayersRepo(),
			r.Actions(),
			r.PlayerProgress(),
			r.Seasons(),
		)
	}
	return r.players
}

func (r *Registry) Effects() *effects.Effects {
	if r.effects == nil {
		r.effects = effects.NewEffects(r.EffectsRepo(), r.InventoriesRepo())
	}
	return r.effects
}

func (r *Registry) Inventories() *inventories.Inventories {
	if r.inventories == nil {
		r.inventories = inventories.NewInventories(r.InventoriesRepo(), r.Effects(), r.ItemsRepo())
	}
	return r.inventories
}

func (r *Registry) Activities() *activities.Activities {
	if r.activities == nil {
		r.activities = activities.NewActivities(r.ActivitiesRepo())
	}
	return r.activities
}

func (r *Registry) ActivityFilters() *activity_filters.ActivityFilters {
	if r.activityFilters == nil {
		r.activityFilters = activity_filters.NewActivityFilters(r.ActivityFiltersRepo())
	}
	return r.activityFilters
}

func (r *Registry) Items() *items.Items {
	if r.items == nil {
		r.items = items.NewItems(r.ItemsRepo())
	}
	return r.items
}

func (r *Registry) Board() *board.Board {
	if r.board == nil {
		r.board = board.NewBoard(r.Actions(), r.PlayerProgress(), r.Cells(), r.Worlds())
	}
	return r.board
}

func (r *Registry) Genres() *genres.Genres {
	if r.genres == nil {
		r.genres = genres.NewGenres(r.GenresRepo())
	}
	return r.genres
}

func (r *Registry) Reviews() *reviews.Reviews {
	if r.reviews == nil {
		r.reviews = reviews.NewReviews(r.ReviewsRepo())
	}
	return r.reviews
}

func (r *Registry) Outboxes() *outboxes.Outboxes {
	if r.outboxes == nil {
		r.outboxes = outboxes.NewOutboxes(r.OutboxesRepo())
	}
	return r.outboxes
}

func (r *Registry) EventStats() *event_stats.EventStats {
	if r.eventStats == nil {
		r.eventStats = event_stats.NewEventStats(r.EventStatsCachedRepo())
	}
	return r.eventStats
}

func (r *Registry) StreamTracker() *stream_tracker.StreamTracker {
	if r.streamTracker == nil {
		r.streamTracker = stream_tracker.NewStreamTracker(r.StreamTrackerRepo(), r.PlayersRepo())
	}
	return r.streamTracker
}
