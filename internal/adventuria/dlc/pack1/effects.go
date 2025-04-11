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
		EffectTypeNothing:                     adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnAny),
		EffectTypePointsIncrement:             adventuria.NewEffect(adventuria.Int, adventuria.EffectUseOnAny),
		EffectTypeJailEscape:                  adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnAny),
		EffectTypeDiceMultiplier:              adventuria.NewEffect(adventuria.Int, adventuria.EffectUseOnRoll),
		EffectTypeDiceIncrement:               adventuria.NewEffect(adventuria.Int, adventuria.EffectUseOnRoll),
		EffectTypeChangeDices:                 adventuria.NewEffectWithSource(adventuria.DiceEffectSourceReceiver, adventuria.EffectUseOnRoll),
		EffectTypeSafeDrop:                    adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnDrop),
		EffectTypeTimerIncrement:              adventuria.NewEffect(adventuria.Int, adventuria.EffectUseOnAny),
		EffectTypeRollReverse:                 adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnRoll),
		EffectTypeDropInventory:               adventuria.NewEffect(adventuria.Bool, adventuria.EffectUseOnAny),
		EffectTypeCellPointsDivide:            adventuria.NewEffect(adventuria.Int, adventuria.EffectUseOnChooseResult),
		EffectTypeTeleportToRandomCellByTypes: adventuria.NewEffectWithSource(adventuria.DefaultEffectSourceReceiver, adventuria.EffectUseOnAny),
	})
}
