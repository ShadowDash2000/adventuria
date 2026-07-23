package adventuria

import (
	"adventuria/internal/adventuria/action_events"
	actionEventsRepo "adventuria/internal/adventuria/action_events/repository"
	"adventuria/internal/adventuria/actions"
	rollWheelRepo "adventuria/internal/adventuria/actions/custom/roll_wheel/repository"
	actionsRepo "adventuria/internal/adventuria/actions/repository"
	"adventuria/internal/adventuria/activities"
	activitiesRepo "adventuria/internal/adventuria/activities/repository"
	"adventuria/internal/adventuria/activity_filters"
	activityFiltersRepo "adventuria/internal/adventuria/activity_filters/repository"
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/cell_events_schedules"
	cellEventsSchedulesRepo "adventuria/internal/adventuria/cell_events_schedules/repository"
	"adventuria/internal/adventuria/cells"
	cellsRepo "adventuria/internal/adventuria/cells/repository"
	"adventuria/internal/adventuria/companies"
	companiesRepo "adventuria/internal/adventuria/companies/repository"
	"adventuria/internal/adventuria/effects"
	effectsRepo "adventuria/internal/adventuria/effects/repository"
	"adventuria/internal/adventuria/event_stats"
	eventStatsRepo "adventuria/internal/adventuria/event_stats/repository"
	"adventuria/internal/adventuria/game_types"
	gameTypesRepo "adventuria/internal/adventuria/game_types/repository"
	"adventuria/internal/adventuria/games/cheapshark"
	cheapSharkRepo "adventuria/internal/adventuria/games/cheapshark/repository"
	"adventuria/internal/adventuria/games/github"
	githubRepo "adventuria/internal/adventuria/games/github/repository"
	"adventuria/internal/adventuria/games/how_long_to_beat"
	hltbRepo "adventuria/internal/adventuria/games/how_long_to_beat/repository"
	"adventuria/internal/adventuria/games/igdb"
	igdbRepo "adventuria/internal/adventuria/games/igdb/repository"
	"adventuria/internal/adventuria/games/steam_spy"
	steamSpyRepo "adventuria/internal/adventuria/games/steam_spy/repository"
	"adventuria/internal/adventuria/genres"
	genresRepo "adventuria/internal/adventuria/genres/repository"
	"adventuria/internal/adventuria/inventories"
	inventoriesRepo "adventuria/internal/adventuria/inventories/repository"
	"adventuria/internal/adventuria/items"
	itemsRepo "adventuria/internal/adventuria/items/repository"
	"adventuria/internal/adventuria/outboxes"
	outboxesRepo "adventuria/internal/adventuria/outboxes/repository"
	"adventuria/internal/adventuria/platforms"
	platformsRepo "adventuria/internal/adventuria/platforms/repository"
	"adventuria/internal/adventuria/player_progress"
	progressRepo "adventuria/internal/adventuria/player_progress/repository"
	"adventuria/internal/adventuria/player_stats"
	playerStatsRepo "adventuria/internal/adventuria/player_stats/repository"
	"adventuria/internal/adventuria/players"
	playersRepo "adventuria/internal/adventuria/players/repository"
	"adventuria/internal/adventuria/reviews"
	reviewsRepo "adventuria/internal/adventuria/reviews/repository"
	"adventuria/internal/adventuria/seasons"
	seasonsRepo "adventuria/internal/adventuria/seasons/repository"
	"adventuria/internal/adventuria/settings"
	settingsRepo "adventuria/internal/adventuria/settings/repository"
	"adventuria/internal/adventuria/stream_tracker"
	streamTrackerRepo "adventuria/internal/adventuria/stream_tracker/repository"
	"adventuria/internal/adventuria/tags"
	tagsRepo "adventuria/internal/adventuria/tags/repository"
	"adventuria/internal/adventuria/themes"
	themesRepo "adventuria/internal/adventuria/themes/repository"
	"adventuria/internal/adventuria/worlds"
	worldsRepo "adventuria/internal/adventuria/worlds/repository"
	"log/slog"
	"os"

	"github.com/pocketbase/pocketbase/core"
)

type Registry struct {
	pb     core.App
	logger *slog.Logger

	// repos
	seasonsRepo             *seasonsRepo.Repository
	settingsRepo            *settingsRepo.Repository
	cellsRepo               *cellsRepo.Repository
	worldsRepo              *worldsRepo.Repository
	actionsRepo             *actionsRepo.Repository
	progressRepo            *progressRepo.Repository
	playersRepo             *playersRepo.Repository
	playerStatsRepo         *playerStatsRepo.Repository
	inventoriesRepo         *inventoriesRepo.Repository
	effectsRepo             *effectsRepo.Repository
	activitiesRepo          *activitiesRepo.Repository
	activityFiltersRepo     *activityFiltersRepo.Repository
	itemsRepo               *itemsRepo.Repository
	genresRepo              *genresRepo.Repository
	reviewsRepo             *reviewsRepo.Repository
	rollWheelRepo           *rollWheelRepo.Repository
	outboxesRepo            *outboxesRepo.Repository
	relationRepo            *activitiesRepo.RelationRepository
	eventStatsRepo          *eventStatsRepo.Repository
	eventStatsCachedRepo    *eventStatsRepo.CachedRepository
	streamTrackerRepo       *streamTrackerRepo.Repository
	platformsRepo           *platformsRepo.Repository
	companiesRepo           *companiesRepo.Repository
	tagsRepo                *tagsRepo.Repository
	themesRepo              *themesRepo.Repository
	gameTypesRepo           *gameTypesRepo.Repository
	hltbRepo                *hltbRepo.Repository
	hltbRemoteRepo          *hltbRepo.RemoteRepository
	steamSpyRepo            *steamSpyRepo.Repository
	steamSpyRemoteRepo      *steamSpyRepo.RemoteRepository
	cheapSharkRepo          *cheapSharkRepo.Repository
	cheapSharkRemoteRepo    *cheapSharkRepo.RemoteRepository
	githubRepo              *githubRepo.Repository
	igdbRepo                *igdbRepo.Repository
	igdbRemoteRepo          *igdbRepo.RemoteRepository
	actionEventsRepo        *actionEventsRepo.Repository
	cellEventsSchedulesRepo *cellEventsSchedulesRepo.Repository

	// services
	seasons             *seasons.Seasons
	settings            *settings.Settings
	worlds              *worlds.Worlds
	cells               *cells.Cells
	actions             *actions.Actions
	progress            *player_progress.PlayerProgress
	players             *players.Players
	playerStats         *player_stats.PlayerStats
	effects             *effects.Effects
	inventories         *inventories.Inventories
	activities          *activities.Activities
	activityFilters     *activity_filters.ActivityFilters
	items               *items.Items
	board               *board.Board
	genres              *genres.Genres
	reviews             *reviews.Reviews
	outboxes            *outboxes.Outboxes
	eventStats          *event_stats.EventStats
	streamTracker       *stream_tracker.StreamTracker
	platforms           *platforms.Platforms
	companies           *companies.Companies
	tags                *tags.Tags
	themes              *themes.Themes
	gameTypes           *game_types.GameTypes
	hltb                *how_long_to_beat.HowLongToBeat
	steamSpy            *steam_spy.SteamSpy
	cheapShark          *cheapshark.CheapShark
	github              *github.Github
	igdb                *igdb.IGDB
	actionEvents        *action_events.ActionEvents
	cellEventsSchedules *cell_events_schedules.CellEventsSchedules
}

func NewRegistry(pb core.App, logger *slog.Logger) *Registry {
	return &Registry{
		pb:     pb,
		logger: logger,
	}
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

func (r *Registry) PlayerStatsRepo() *playerStatsRepo.Repository {
	if r.playerStatsRepo == nil {
		r.playerStatsRepo = playerStatsRepo.NewRepository(r.pb)
	}
	return r.playerStatsRepo
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

func (r *Registry) PlatformsRepo() *platformsRepo.Repository {
	if r.platformsRepo == nil {
		r.platformsRepo = platformsRepo.NewRepository(r.pb)
	}
	return r.platformsRepo
}

func (r *Registry) CompaniesRepo() *companiesRepo.Repository {
	if r.companiesRepo == nil {
		r.companiesRepo = companiesRepo.NewRepository(r.pb)
	}
	return r.companiesRepo
}

func (r *Registry) TagsRepo() *tagsRepo.Repository {
	if r.tagsRepo == nil {
		r.tagsRepo = tagsRepo.NewRepository(r.pb)
	}
	return r.tagsRepo
}

func (r *Registry) ThemesRepo() *themesRepo.Repository {
	if r.themesRepo == nil {
		r.themesRepo = themesRepo.NewRepository(r.pb)
	}
	return r.themesRepo
}

func (r *Registry) GameTypesRepo() *gameTypesRepo.Repository {
	if r.gameTypesRepo == nil {
		r.gameTypesRepo = gameTypesRepo.NewRepository(r.pb)
	}
	return r.gameTypesRepo
}

func (r *Registry) HltbRepo() *hltbRepo.Repository {
	if r.hltbRepo == nil {
		r.hltbRepo = hltbRepo.NewRepository(r.pb)
	}
	return r.hltbRepo
}

func (r *Registry) HltbRemoteRepo() *hltbRepo.RemoteRepository {
	if r.hltbRemoteRepo == nil {
		r.hltbRemoteRepo = hltbRepo.NewRemoteRepository(r.GithubRepo())
	}
	return r.hltbRemoteRepo
}

func (r *Registry) SteamSpyRepo() *steamSpyRepo.Repository {
	if r.steamSpyRepo == nil {
		r.steamSpyRepo = steamSpyRepo.NewRepository(r.pb)
	}
	return r.steamSpyRepo
}

func (r *Registry) SteamSpyRemoteRepo() *steamSpyRepo.RemoteRepository {
	if r.steamSpyRemoteRepo == nil {
		r.steamSpyRemoteRepo = steamSpyRepo.NewRemoteRepository(r.GithubRepo())
	}
	return r.steamSpyRemoteRepo
}

func (r *Registry) CheapSharkRepo() *cheapSharkRepo.Repository {
	if r.cheapSharkRepo == nil {
		r.cheapSharkRepo = cheapSharkRepo.NewRepository(r.pb)
	}
	return r.cheapSharkRepo
}

func (r *Registry) CheapSharkRemoteRepo() *cheapSharkRepo.RemoteRepository {
	if r.cheapSharkRemoteRepo == nil {
		r.cheapSharkRemoteRepo = cheapSharkRepo.NewRemoteRepository(r.GithubRepo())
	}
	return r.cheapSharkRemoteRepo
}

func (r *Registry) GithubRepo() *githubRepo.Repository {
	if r.githubRepo == nil {
		r.githubRepo = githubRepo.NewRepository()
	}
	return r.githubRepo
}

func (r *Registry) IGDBRepo() *igdbRepo.Repository {
	if r.igdbRepo == nil {
		r.igdbRepo = igdbRepo.NewRepository(r.pb)
	}
	return r.igdbRepo
}

func (r *Registry) IGDBRemoteRepo() *igdbRepo.RemoteRepository {
	twitchClientId, _ := os.LookupEnv("TWITCH_CLIENT_ID")
	twitchClientSecret, _ := os.LookupEnv("TWITCH_CLIENT_SECRET")

	if r.igdbRemoteRepo == nil {
		r.igdbRemoteRepo = igdbRepo.NewRemoteRepository(twitchClientId, twitchClientSecret)
	}
	return r.igdbRemoteRepo
}

func (r *Registry) ActionEventsRepo() *actionEventsRepo.Repository {
	if r.actionEventsRepo == nil {
		r.actionEventsRepo = actionEventsRepo.NewRepository(r.pb)
	}
	return r.actionEventsRepo
}

func (r *Registry) CellEventsSchedulesRepo() *cellEventsSchedulesRepo.Repository {
	if r.cellEventsSchedulesRepo == nil {
		r.cellEventsSchedulesRepo = cellEventsSchedulesRepo.NewRepository(r.pb)
	}
	return r.cellEventsSchedulesRepo
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
		r.actions = actions.NewActions(r.ActionsRepo(), r.Worlds(), r.Cells(), r.ActionEvents())
	}
	return r.actions
}

func (r *Registry) PlayerProgress() *player_progress.PlayerProgress {
	if r.progress == nil {
		r.progress = player_progress.NewPlayerProgress(r.PlayerProgressRepo(), r.PlayerProgressRepo(), r.Worlds())
	}
	return r.progress
}

func (r *Registry) Players() *players.Players {
	if r.players == nil {
		r.players = players.NewPlayers(
			r.PlayersRepo(),
			r.Actions(),
			r.PlayerProgress(),
			r.PlayerStats(),
			r.Seasons(),
		)
	}
	return r.players
}

func (r *Registry) PlayerStats() *player_stats.PlayerStats {
	if r.playerStats == nil {
		r.playerStats = player_stats.NewPlayerStats(r.PlayerStatsRepo())
	}
	return r.playerStats
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
		r.streamTracker = stream_tracker.NewStreamTracker(r.logger, r.StreamTrackerRepo(), r.PlayersRepo())
	}
	return r.streamTracker
}

func (r *Registry) Platforms() *platforms.Platforms {
	if r.platforms == nil {
		r.platforms = platforms.NewPlatforms(r.PlatformsRepo())
	}
	return r.platforms
}

func (r *Registry) Companies() *companies.Companies {
	if r.companies == nil {
		r.companies = companies.NewCompanies(r.CompaniesRepo())
	}
	return r.companies
}

func (r *Registry) Tags() *tags.Tags {
	if r.tags == nil {
		r.tags = tags.NewTags(r.TagsRepo())
	}
	return r.tags
}

func (r *Registry) Themes() *themes.Themes {
	if r.themes == nil {
		r.themes = themes.NewThemes(r.ThemesRepo())
	}
	return r.themes
}

func (r *Registry) GameTypes() *game_types.GameTypes {
	if r.gameTypes == nil {
		r.gameTypes = game_types.NewGameTypes(r.GameTypesRepo())
	}
	return r.gameTypes
}

func (r *Registry) HLTB() *how_long_to_beat.HowLongToBeat {
	if r.hltb == nil {
		r.hltb = how_long_to_beat.NewHowLongToBeat(r.HltbRepo(), r.HltbRemoteRepo())
	}
	return r.hltb
}

func (r *Registry) SteamSpy() *steam_spy.SteamSpy {
	if r.steamSpy == nil {
		r.steamSpy = steam_spy.NewSteamSpy(r.SteamSpyRepo(), r.SteamSpyRemoteRepo())
	}
	return r.steamSpy
}

func (r *Registry) CheapShark() *cheapshark.CheapShark {
	if r.cheapShark == nil {
		r.cheapShark = cheapshark.NewCheapShark(r.CheapSharkRepo(), r.CheapSharkRemoteRepo())
	}
	return r.cheapShark
}

func (r *Registry) IGDB() *igdb.IGDB {
	if r.igdb == nil {
		r.igdb = igdb.NewIGDB(
			r.IGDBRepo(),
			r.IGDBRemoteRepo(),
			r.Activities(),
			r.Platforms(),
			r.Companies(),
			r.Tags(),
			r.Themes(),
			r.Genres(),
			r.GameTypes(),
			r.HLTB(),
			r.SteamSpy(),
			r.CheapShark(),
			r.Settings(),
		)
	}
	return r.igdb
}

func (r *Registry) ActionEvents() *action_events.ActionEvents {
	if r.actionEvents == nil {
		r.actionEvents = action_events.NewActionEvents(r.ActionEventsRepo())
	}
	return r.actionEvents
}

func (r *Registry) CellEventsSchedules() *cell_events_schedules.CellEventsSchedules {
	if r.cellEventsSchedules == nil {
		r.cellEventsSchedules = cell_events_schedules.NewCellEventsSchedules(
			r.CellEventsSchedulesRepo(),
			r.Cells(),
			r.Effects(),
			r.ActionEvents(),
			r.Players(),
			r.Settings(),
		)
	}
	return r.cellEventsSchedules
}
