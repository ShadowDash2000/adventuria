package adventuria

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
	RegisterEffects(map[string]EffectCreator{
		EffectTypeNothing:                     NewEffectInt(),
		EffectTypePointsIncrement:             NewEffectInt(),
		EffectTypeJailEscape:                  NewEffectBool(),
		EffectTypeDiceMultiplier:              NewEffectInt(),
		EffectTypeDiceIncrement:               NewEffectInt(),
		EffectTypeChangeDices:                 NewEffectSliceWithSource(Dices),
		EffectTypeSafeDrop:                    NewEffectBool(),
		EffectTypeTimerIncrement:              NewEffectInt(),
		EffectTypeRollReverse:                 NewEffectBool(),
		EffectTypeDropInventory:               NewEffectBool(),
		EffectTypeCellPointsDivide:            NewEffectInt(),
		EffectTypeTeleportToRandomCellByTypes: NewEffectSlice(),
	})
}
