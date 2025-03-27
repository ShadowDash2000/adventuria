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
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		EffectTypeNothing:                     adventuria.NewEffect(adventuria.Bool),
		EffectTypePointsIncrement:             adventuria.NewEffect(adventuria.Int),
		EffectTypeJailEscape:                  adventuria.NewEffect(adventuria.Bool),
		EffectTypeDiceMultiplier:              adventuria.NewEffect(adventuria.Int),
		EffectTypeDiceIncrement:               adventuria.NewEffect(adventuria.Int),
		EffectTypeChangeDices:                 adventuria.NewEffectWithSource(adventuria.DiceEffectSourceReceiver),
		EffectTypeSafeDrop:                    adventuria.NewEffect(adventuria.Bool),
		EffectTypeTimerIncrement:              adventuria.NewEffect(adventuria.Int),
		EffectTypeRollReverse:                 adventuria.NewEffect(adventuria.Bool),
		EffectTypeDropInventory:               adventuria.NewEffect(adventuria.Bool),
		EffectTypeCellPointsDivide:            adventuria.NewEffect(adventuria.Int),
		EffectTypeTeleportToRandomCellByTypes: adventuria.NewEffect(adventuria.Slice),
	})
}
