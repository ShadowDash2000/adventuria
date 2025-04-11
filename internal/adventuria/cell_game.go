package adventuria

type CellGame struct {
	CellWheelBase
}

func NewCellGame() CellCreator {
	return func() Cell {
		return &CellGame{}
	}
}

func (c *CellGame) NextStep(user *User) string {
	nextStepType := ""

	switch user.LastAction.Type() {
	case ActionTypeRoll,
		ActionTypeReroll:
		nextStepType = ActionTypeChooseGame
	case ActionTypeChooseGame:
		nextStepType = ActionTypeDone
	case ActionTypeDone,
		ActionTypeDrop:
		nextStepType = ActionTypeRoll
	default:
		nextStepType = ActionTypeChooseGame
	}

	return nextStepType
}

func (c *CellGame) Roll(_ *User) (*WheelRollResult, error) {
	res := &WheelRollResult{}

	return res, nil
}
