package adventuria

import (
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

// EffectVerifier
// Binds hooks on effect's collection for record creation and update
// that verifies that an effect type really exists and calls Verify()
// method of an effect that should try to parse record's value
type EffectVerifier struct{}

func NewEffectVerifier() *EffectVerifier {
	ef := &EffectVerifier{}
	ef.bindHooks()
	return ef
}

func (ef *EffectVerifier) bindHooks() {
	PocketBase.OnRecordCreate(CollectionEffects).BindFunc(func(e *core.RecordEvent) error {
		if err := ef.Verify(e.Record.GetString("type"), e.Record.GetString("value")); err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordUpdate(CollectionEffects).BindFunc(func(e *core.RecordEvent) error {
		if err := ef.Verify(e.Record.GetString("type"), e.Record.GetString("value")); err != nil {
			return err
		}
		return e.Next()
	})
}

func (ef *EffectVerifier) Verify(effectType, value string) error {
	effectCreator, ok := effectsList[effectType]
	if !ok {
		return errors.New("unknown effect type")
	}

	effect := effectCreator(nil, core.NewRecord(GameCollections.Get(CollectionEffects)))

	return effect.Verify(value)
}
