package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

type AddGameTagEffect struct {
	adventuria.EffectBase
}

func (ef *AddGameTagEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if ctx.InvItemID == e.ItemId {
				if !adventuria.GameActions.CanDo(ctx.User, "rollWheel") {
					return errors.New("addGameTag: can't add tag while rollWheel isn't available")
				}

				if tagID, ok := e.Request["tag_id"].(string); ok {
					_, err := ef.fetchGameTagByID(tagID)
					if err != nil {
						return err
					}

					filter := ctx.User.LastAction().CustomGameFilter()
					filter.Tags = append(filter.Tags, tagID)

					callback()
				} else {
					return errors.New("addGameTag: tag_id not found")
				}
			}

			return e.Next()
		}),
	}
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
