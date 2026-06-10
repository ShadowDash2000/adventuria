package effects

import "adventuria/internal/adventuria_new/model"

func effectInfoToEffect(info *model.EffectInfo) (model.Effect, error) {
	return Create(*info)
}

func effectInfosToEffects(infos []*model.EffectInfo) ([]model.Effect, error) {
	effects := make([]model.Effect, len(infos))
	for i, info := range infos {
		effect, err := effectInfoToEffect(info)
		if err != nil {
			return nil, err
		}
		effects[i] = effect
	}
	return effects, nil
}
