package adventuria

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
)

type CellBase struct {
	core.BaseRecordProxy
}

func NewCellFromRecord(record *core.Record) (Cell, error) {
	t := CellType(record.GetString("type"))

	cellCreator, ok := CellsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown cell type: %s", t)
	}

	cell := cellCreator()
	cell.SetProxyRecord(record)

	return cell, nil
}

func (c *CellBase) ID() string {
	return c.Id
}

func (c *CellBase) IsActive() bool {
	return c.GetBool("isActive")
}

func (c *CellBase) Sort() int {
	return c.GetInt("sort")
}

func (c *CellBase) Type() CellType {
	return CellType(c.GetString("type"))
}

func (c *CellBase) SetType(t CellType) {
	c.Set("type", t)
}

func (c *CellBase) Preset() string {
	return c.GetString("preset")
}

func (c *CellBase) AudioPresets() []string {
	return c.GetStringSlice("audioPresets")
}

func (c *CellBase) Icon() string {
	return c.GetString("icon")
}

func (c *CellBase) Name() string {
	return c.GetString("name")
}

func (c *CellBase) Points() int {
	return c.GetInt("points")
}

func (c *CellBase) Description() string {
	return c.GetString("description")
}

func (c *CellBase) Color() string {
	return c.GetString("color")
}

func (c *CellBase) CantDrop() bool {
	return c.GetBool("cantDrop")
}

func (c *CellBase) CantReroll() bool {
	return c.GetBool("cantReroll")
}

func (c *CellBase) CantChooseAfterDrop() bool {
	return c.GetBool("cantChooseAfterDrop")
}

func (c *CellBase) IsSafeDrop() bool {
	return c.GetBool("isSafeDrop")
}

func (c *CellBase) NextStep(_ *User) string {
	return ActionTypeRollDice
}

func (c *CellBase) OnCellReached(_ *User) error {
	return nil
}

type CellTypeSourceGiver struct {
	source []string
}

func NewCellTypeSourceGiver(source []string) EffectSourceGiver[CellType] {
	return &CellTypeSourceGiver{source: source}
}

func (cg *CellTypeSourceGiver) Slice() []CellType {
	var res []CellType
	for _, cellType := range cg.source {
		res = append(res, CellType(cellType))
	}
	return res
}
