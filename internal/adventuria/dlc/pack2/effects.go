package pack2

import "adventuria/internal/adventuria"

const (
	EffectTypeChangeNextStepType        = "changeNextStepType"
	EffectTypeCustomGame                = "customGame"
	EffectTypeGoToJail                  = "goToJail"
	EffectTypeTeleportToRandomCellByIds = "teleportToRandomCellByIds"
)

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		EffectTypeChangeNextStepType:        adventuria.NewEffect(adventuria.String, adventuria.EffectUseOnAny),
		EffectTypeGoToJail:                  adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnAny),
		EffectTypeTeleportToRandomCellByIds: adventuria.NewEffect(adventuria.Slice, adventuria.EffectUseOnAny),
	})
}
