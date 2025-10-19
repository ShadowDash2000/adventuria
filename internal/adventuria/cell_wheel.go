package adventuria

type CellWheel interface {
	Cell
	Roll(User) (*WheelRollResult, error)
	GetItems(User) ([]*WheelItem, error)
}

type WheelRollResult struct {
	FillerItems []*WheelItem `json:"fillerItems"`
	WinnerId    string       `json:"winnerId"`
}

type WheelItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type CellWheelBase struct {
	CellBase
}

func (c *CellWheelBase) Roll(_ User) (*WheelRollResult, error) {
	panic("implement me")
}

func (c *CellWheelBase) GetItems(_ User) ([]*WheelItem, error) {
	panic("implement me")
}
