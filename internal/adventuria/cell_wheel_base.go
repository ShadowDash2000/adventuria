package adventuria

import "github.com/pocketbase/pocketbase/core"

type WheelRollResult struct {
	FillerItems []*WheelItem     `json:"fillerItems"`
	WinnerId    string           `json:"winnerId"`
	Collection  *core.Collection `json:"collection"`
	EffectUse   EffectUse        `json:"effectUse"`
}

type WheelItem struct {
	Name string
	Icon string
}

type CellWheel interface {
	Cell
	Roll(*User) (*WheelRollResult, error)
}

type CellWheelBase struct {
	CellBase
}
