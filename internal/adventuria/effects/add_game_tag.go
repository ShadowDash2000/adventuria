package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"errors"
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
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if ctx.InvItemID == e.InvItemId {
				cell, ok := ctx.User.CurrentCell()
				if !ok {
					return errors.New("addGameTag: current cell not found")
				}

				cellGame, ok := cell.(*cells.CellGame)
				if !ok {
					return errors.New("addGameTag: current cell isn't game cell")
				}

				if tagID, ok := e.Request["tag_id"].(string); ok {
					_, err := ef.fetchGameTagByID(tagID)
					if err != nil {
						return err
					}

					filter := ctx.User.LastAction().CustomGameFilter()
					if index := slices.Index(filter.Tags, tagID); index != -1 {
						return errors.New("addGameTag: tag already exists")
					}

					filter.Tags = append(filter.Tags, tagID)
					if err := cellGame.CheckCustomFilter(ctx.User); err != nil {
						return fmt.Errorf("addGameTag: %w", err)
					}

					callback()
				} else {
					return errors.New("addGameTag: tag_id not found")
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
