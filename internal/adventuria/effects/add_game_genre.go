package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"fmt"
	"slices"
)

type AddGameGenreEffect struct {
	adventuria.EffectRecord
}

func (ef *AddGameGenreEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if !adventuria.GameActions.CanDo(appCtx, ctx.User, "rollWheel") {
		return false
	}

	cell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if cell.Type() != "game" {
		return false
	}

	if filterId := cell.Filter(); filterId != "" {
		filterRecord, err := appCtx.App.FindRecordById(
			schema.CollectionActivityFilter,
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
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				cell, ok := ctx.User.CurrentCell()
				if !ok {
					return result.Err("internal error: current cell not found"),
						errors.New("addGameGenre: current cell not found")
				}

				cellGame, ok := cell.(adventuria.CellWheel)
				if !ok {
					return result.Err("current cell isn't wheel cell"), nil
				}

				if genreId, ok := e.Data["genre_id"].(string); ok {
					_, err := e.AppContext.App.FindRecordById(
						schema.CollectionGenres,
						genreId,
					)
					if err != nil {
						return result.Err("genre_id not found"), nil
					}

					filter, err := ctx.User.LastAction().CustomActivityFilter()
					if err != nil {
						return result.Err("internal error: can't get custom activity filter"),
							fmt.Errorf("addGameGenre: %w", err)
					}

					if index := slices.Index(filter.Tags, genreId); index != -1 {
						return result.Err("genre already exists"), nil
					}

					filter.Genres = append(filter.Genres, genreId)
					if err = cellGame.RefreshItems(e.AppContext, ctx.User); err != nil {
						return result.Err("internal error: can't refresh cell items"),
							fmt.Errorf("addGameGenre: %w", err)
					}

					ctx.User.LastAction().SetCustomActivityFilter(*filter)

					callback(e.AppContext)
				} else {
					return result.Err("genre_id not specified"), nil
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *AddGameGenreEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *AddGameGenreEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}
