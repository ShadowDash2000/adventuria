package adventuria

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

// ensure CellBase implements Cell
var _ Cell = (*CellBase)(nil)

type CellBase struct {
	core.BaseRecordProxy
}

func NewCellFromRecord(record *core.Record) (Cell, error) {
	t := CellType(record.GetString("type"))

	cellCreator, ok := cellsList[t]
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

func (c *CellBase) Sort() int {
	return c.GetInt("sort")
}

func (c *CellBase) Type() CellType {
	return CellType(c.GetString("type"))
}

func (c *CellBase) SetType(t CellType) {
	c.Set("type", t)
}

func (c *CellBase) Filter() string {
	return c.GetString("filter")
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

func (c *CellBase) Coins() int {
	return c.GetInt("coins")
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

func (c *CellBase) IsSafeDrop() bool {
	return c.GetBool("isSafeDrop")
}

func (c *CellBase) OnCellReached(_ *CellReachedContext) error {
	return nil
}

func (c *CellBase) Verify(_ string) error {
	panic("implement me")
}

func (c *CellBase) DecodeValue(_ string) (any, error) {
	panic("implement me")
}

func (c *CellBase) Value() string {
	return c.GetString("value")
}
