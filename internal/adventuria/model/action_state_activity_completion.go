package model

type ActionActivityCompletionState struct {
	Points                  int
	DoneEnergyConsume       int
	DropEnergyConsume       int
	Coins                   int
	CanDrop                 bool
	CanReroll               bool
	IsSafeDrop              bool
	IsCustomFilterAllowed   bool
	IsChangeActivityAllowed bool
	GiveItemWheel           bool
}

func (a *ActionActivityCompletionState) Clone() *ActionActivityCompletionState {
	if a == nil {
		return nil
	}

	return &ActionActivityCompletionState{
		Points:                  a.Points,
		DoneEnergyConsume:       a.DoneEnergyConsume,
		DropEnergyConsume:       a.DropEnergyConsume,
		Coins:                   a.Coins,
		CanDrop:                 a.CanDrop,
		CanReroll:               a.CanReroll,
		IsSafeDrop:              a.IsSafeDrop,
		IsCustomFilterAllowed:   a.IsCustomFilterAllowed,
		IsChangeActivityAllowed: a.IsChangeActivityAllowed,
		GiveItemWheel:           a.GiveItemWheel,
	}
}
