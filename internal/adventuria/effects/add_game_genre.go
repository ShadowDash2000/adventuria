package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
	"slices"
)

type AddGameGenreEffect struct {
	adventuria.EffectBase
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

				if cell.Type() != "game" {
					return &event.Result{
						Success: false,
						Error:   "current cell isn't game cell",
					}, nil
				}

				if filterId := cell.Filter(); filterId != "" {
					filterRecord, err := adventuria.PocketBase.FindRecordById(
						adventuria.CollectionActivityFilter,
						filterId,
					)
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "filter not found",
						}, fmt.Errorf("addGameGenre: %w", err)
					}

					if len(filterRecord.GetStringSlice("activities")) > 0 {
						return &event.Result{
							Success: false,
							Error:   "current cell filter has activities, can't add genre",
						}, nil
					}
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

func (ef *AddGameGenreEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}
