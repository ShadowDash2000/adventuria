package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
	"slices"
)

type AddGameGenreEffect struct {
	adventuria.EffectRecord
}

func (ef *AddGameGenreEffect) CanUse(ctx adventuria.EffectContext) bool {
	if !adventuria.GameActions.CanDo(ctx.User, "rollWheel") {
		return false
	}

	cell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	_, ok = cell.(*cells.CellGame)
	if !ok {
		return false
	}

	if cell.Type() != "game" {
		return false
	}

	if filterId := cell.Filter(); filterId != "" {
		filterRecord, err := adventuria.PocketBase.FindRecordById(
			adventuria.CollectionActivityFilter,
			filterId,
		)
		if err != nil {
			return false
		}

		if len(filterRecord.GetStringSlice("developers")) > 0 {
			return false
		}
		if len(filterRecord.GetStringSlice("publishers")) > 0 {
			return false
		}
		if len(filterRecord.GetStringSlice("activities")) > 0 {
			return false
		}
	}

	return true
}

func (ef *AddGameGenreEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				cell, ok := ctx.User.CurrentCell()
				if !ok {
					return &event.Result{
						Success: false,
						Error:   "current cell not found",
					}, nil
				}

				cellGame, ok := cell.(*cells.CellGame)
				if !ok {
					return &event.Result{
						Success: false,
						Error:   "current cell isn't game cell",
					}, nil
				}

				if genreId, ok := e.Request["genre_id"].(string); ok {
					_, err := adventuria.PocketBase.FindRecordById(
						adventuria.CollectionGenres,
						genreId,
					)
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "genre_id not found",
						}, fmt.Errorf("addGameGenre: %w", err)
					}

					filter := ctx.User.LastAction().CustomActivityFilter()
					if index := slices.Index(filter.Tags, genreId); index != -1 {
						return &event.Result{
							Success: false,
							Error:   "genre already exists",
						}, nil
					}

					filter.Genres = append(filter.Genres, genreId)
					if err = cellGame.RefreshItems(ctx.User); err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't refresh cell items",
						}, fmt.Errorf("addGameGenre: %w", err)
					}

					callback()
				} else {
					return &event.Result{
						Success: false,
						Error:   "genre_id not found",
					}, nil
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *AddGameGenreEffect) Verify(_ string) error {
	return nil
}

func (ef *AddGameGenreEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}
