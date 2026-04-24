package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
)

type WorldEffects struct {
	player     Player
	unsubGroup event.UnsubGroup
}

func NewWorldEffects(player Player) *WorldEffects {
	return &WorldEffects{player: player}
}

func (we *WorldEffects) Subscribe(ctx AppContext, worldId string) error {
	we.unsubGroup.Unsubscribe()

	if worldId == "" {
		return nil
	}

	worldRecord, err := ctx.App.FindRecordById(schema.CollectionsWorlds, worldId)
	if err != nil {
		return err
	}

	errs := ctx.App.ExpandRecord(worldRecord, []string{schema.WorldsSchema.Effects}, nil)
	if len(errs) > 0 {
		return errs[schema.WorldsSchema.Effects]
	}

	effectRecords := worldRecord.ExpandedAll(schema.WorldsSchema.Effects)

	for _, record := range effectRecords {
		effect, err := NewEffectFromRecord(record)
		if err != nil {
			continue
		}

		unsubs, err := effect.Subscribe(EffectContext{
			Player:   we.player,
			Priority: 100,
		}, func(ctx AppContext) {})

		if err == nil {
			we.unsubGroup.Add(unsubs...)
		}
	}

	return nil
}

func (we *WorldEffects) Close() {
	we.unsubGroup.Unsubscribe()
}
