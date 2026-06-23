package effects

import (
	"adventuria/internal/adventuria/schema"
	repo "adventuria/internal/adventuria_new/effects/repository"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func BindHooks(pb core.App) {
	pb.OnRecordValidate(schema.CollectionEffects).BindFunc(func(e *core.RecordEvent) error {
		err := verify(e.Context, repo.RecordToEffect(e.Record))
		if err != nil {
			return err
		}
		return e.Next()
	})
}

func verify(ctx context.Context, effectInfo *model.EffectInfo) error {
	effectValue := effectInfo.Value()

	effectDef, ok := Get(effectInfo.Type())
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownEffectType, effectInfo.Type())
	}

	effect := effectDef.new(*effectInfo)
	effectVerifiable, ok := effect.(model.Verifiable)
	if !ok {
		if effectValue == "" {
			return nil
		}
		return errors.New("effect type is not verifiable")
	}

	return effectVerifiable.Verify(ctx, effectValue)
}
