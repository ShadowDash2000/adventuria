package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
	"slices"

	"github.com/pocketbase/pocketbase/core"
)

type AddGameTagEffect struct {
	adventuria.EffectBase
}

func (ef *AddGameTagEffect) Subscribe(
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

				if tagID, ok := e.Request["tag_id"].(string); ok {
					_, err := ef.fetchGameTagByID(tagID)
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "tag_id not found",
						}, fmt.Errorf("addGameTag: %w", err)
					}

					filter := ctx.User.LastAction().CustomGameFilter()
					if index := slices.Index(filter.Tags, tagID); index != -1 {
						return &event.Result{
							Success: false,
							Error:   "tag already exists",
						}, nil
					}

					filter.Tags = append(filter.Tags, tagID)
					if err = cellGame.RefreshItems(ctx.User); err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't refresh cell items",
						}, fmt.Errorf("addGameTag: %w", err)
					}

					callback()
				} else {
					return &event.Result{
						Success: false,
						Error:   "tag_id not found",
					}, nil
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *AddGameTagEffect) fetchGameTagByID(tagID string) (*core.Record, error) {
	return adventuria.PocketBase.FindRecordById(
		adventuria.GameCollections.Get(adventuria.CollectionTags),
		tagID,
	)
}

func (ef *AddGameTagEffect) Verify(_ string) error {
	return nil
}

func (ef *AddGameTagEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}
