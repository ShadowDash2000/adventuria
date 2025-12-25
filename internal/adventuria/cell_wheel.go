package adventuria

type CellWheel interface {
	Cell
	Roll(User, RollWheelRequest) (*WheelRollResult, error)
	RefreshItems(User) error
}

type RollWheelRequest map[string]any

type WheelRollResult struct {
	FillerItems []WheelItem `json:"fillerItems"`
	WinnerId    string      `json:"winnerId"`
	Success     bool        `json:"success"`
	Error       string      `json:"error,omitempty"`
}

type WheelItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// ensure CellWheelBase implements CellWheel
var _ CellWheel = (*CellWheelBase)(nil)

type CellWheelBase struct {
	CellBase
}

func (c *CellWheelBase) Roll(_ User, _ RollWheelRequest) (*WheelRollResult, error) {
	panic("implement me")
}

func (c *CellWheelBase) RefreshItems(_ User) error {
	panic("implement me")
}
