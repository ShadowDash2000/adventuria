package pack1

import "adventuria/internal/adventuria"

const (
	EffectTypeNothing                     = "nothing"
	EffectTypePointsIncrement             = "pointsIncrement"
	EffectTypeJailEscape                  = "jailEscape"
	EffectTypeDiceMultiplier              = "diceMultiplier"
	EffectTypeDiceIncrement               = "diceIncrement"
	EffectTypeChangeDices                 = "changeDices"
	EffectTypeSafeDrop                    = "isSafeDrop"
	EffectTypeTimerIncrement              = "timerIncrement"
	EffectTypeRollReverse                 = "rollReverse"
	EffectTypeDropInventory               = "dropInventory"
	EffectTypeCellPointsDivide            = "cellPointsDivide"
	EffectTypeTeleportToRandomCellByTypes = "teleportToRandomCellByTypes"
)

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.Effect{
		EffectTypeNothing:                     adventuria.NewEffectInt(),
		EffectTypePointsIncrement:             adventuria.NewEffectInt(),
		EffectTypeJailEscape:                  adventuria.NewEffectBool(),
		EffectTypeDiceMultiplier:              adventuria.NewEffectInt(),
		EffectTypeDiceIncrement:               adventuria.NewEffectInt(),
		EffectTypeChangeDices:                 adventuria.NewEffectSliceWithSource(adventuria.Dices),
		EffectTypeSafeDrop:                    adventuria.NewEffectBool(),
		EffectTypeTimerIncrement:              adventuria.NewEffectInt(),
		EffectTypeRollReverse:                 adventuria.NewEffectBool(),
		EffectTypeDropInventory:               adventuria.NewEffectBool(),
		EffectTypeCellPointsDivide:            adventuria.NewEffectInt(),
		EffectTypeTeleportToRandomCellByTypes: adventuria.NewEffectSlice(),
	})
}
