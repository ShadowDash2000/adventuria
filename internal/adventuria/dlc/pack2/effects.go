package pack2

import "adventuria/internal/adventuria"

const (
	EffectTypeChangeNextStepType        = "changeNextStepType"
	EffectTypeGoToJail                  = "goToJail"
	EffectTypeTeleportToRandomCellByIds = "teleportToRandomCellByIds"
	EffectChangeCellById                = "changeCellById"
)

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		EffectTypeChangeNextStepType:        adventuria.NewEffect(adventuria.String, adventuria.EffectUseOnAny),
		EffectTypeGoToJail:                  adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnAny),
		EffectTypeTeleportToRandomCellByIds: adventuria.NewEffect(adventuria.Slice, adventuria.EffectUseOnAny),
		EffectChangeCellById:                adventuria.NewEffect(adventuria.String, adventuria.EffectUseOnAny),
	})
}
