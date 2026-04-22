package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

// EffectVerifier
// Binds hooks on effect's collection for record creation and update
// that verifies that an effect type really exists and calls Verify()
// method of an effect that should try to parse record's value
type EffectVerifier struct{}

func NewEffectVerifier(ctx AppContext) *EffectVerifier {
	ef := &EffectVerifier{}
	ef.bindHooks(ctx)
	return ef
}

func (ef *EffectVerifier) bindHooks(ctx AppContext) {
	ctx.App.OnRecordValidate(schema.CollectionEffects).BindFunc(func(e *core.RecordEvent) error {
		if err := ef.Verify(AppContext{App: e.App}, e.Record); err != nil {
			return err
		}
		return e.Next()
	})
}

func (ef *EffectVerifier) Verify(ctx AppContext, record *core.Record) error {
	effectType := record.GetString(schema.EffectSchema.Type)
	effectValue := record.GetString(schema.EffectSchema.Value)

	effectCreator, ok := effectsList[effectType]
	if !ok {
		return errors.New("unknown effect type")
	}

	effect := effectCreator(core.NewRecord(GameCollections.Get(schema.CollectionEffects)))
	effectVerifiable, ok := effect.(EffectVerifiable)

	if !ok {
		if effectValue == "" {
			return nil
		}
		return errors.New("effect type is not verifiable")
	}

	return effectVerifiable.Verify(ctx, effectValue)
}
