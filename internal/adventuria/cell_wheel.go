package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type CellWheel interface {
	Cell
	Roll(User) (*WheelRollResult, error)
}

type WheelRollResult struct {
	FillerItems []*WheelItem     `json:"fillerItems"`
	WinnerId    string           `json:"winnerId"`
	Collection  *core.Collection `json:"collection"`
}

type WheelItem struct {
	Name string
	Icon string
}
