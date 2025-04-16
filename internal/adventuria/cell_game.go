package adventuria

type CellGame struct {
	CellBase
}

func NewCellGame() CellCreator {
	return func() Cell {
		return &CellGame{}
	}
}

func (c *CellGame) NextStep(user *User) string {
	nextStepType := ""

	switch user.LastAction.Type() {
	case ActionTypeRollDice,
		ActionTypeReroll:
		nextStepType = ActionTypeRollWheel
	case ActionTypeRollWheel:
		nextStepType = ActionTypeChooseResult
	case ActionTypeChooseResult,
		ActionTypeDrop:
		nextStepType = ActionTypeRollDice
	default:
		nextStepType = ActionTypeRollWheel
	}

	return nextStepType
}

func (c *CellGame) Roll(_ *User) (*WheelRollResult, error) {
	res := &WheelRollResult{}

	// TODO
	panic("implement me")

	return res, nil
}
