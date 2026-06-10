package custom

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/activity_filters"
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/effects/custom/add_game_genre"
	"adventuria/internal/adventuria_new/genres"
)

func RegisterEffects(
	actions *actions.Actions,
	cells *cells.Cells,
	genres *genres.Genres,
	activityFilters *activity_filters.ActivityFilters,
) {
	effects.Register(
		add_game_genre.NewAddGameGenreEffectDef(
			actions,
			cells,
			genres,
			activityFilters,
		),
	)
}
