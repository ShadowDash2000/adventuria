package adventuria

import "github.com/pocketbase/pocketbase/core"

const (
	CellTypeGame        = "game"
	CellTypeStart       = "start"
	CellTypeJail        = "jail"
	CellTypePreset      = "preset"
	CellTypeItem        = "item"
	CellTypeWheelPreset = "wheelPreset"
)

type Cell struct {
	core.BaseRecordProxy
}

func NewCell(record *core.Record) *Cell {
	c := &Cell{}
	c.SetProxyRecord(record)
	return c
}

func (c *Cell) Sort() int {
	return c.GetInt("sort")
}

func (c *Cell) Type() string {
	return c.GetString("type")
}

func (c *Cell) Preset() string {
	return c.GetString("preset")
}

func (c *Cell) AudioPresets() []string {
	return c.GetStringSlice("audioPresets")
}

func (c *Cell) Icon() string {
	return c.GetString("icon")
}

func (c *Cell) Name() string {
	return c.GetString("name")
}

func (c *Cell) Points() int {
	return c.GetInt("points")
}

func (c *Cell) Description() string {
	return c.GetString("description")
}

func (c *Cell) Color() string {
	return c.GetString("color")
}

func (c *Cell) CantDrop() bool {
	return c.GetBool("cantDrop")
}

func (c *Cell) CantReroll() bool {
	return c.GetBool("cantReroll")
}

func (c *Cell) CantChooseAfterDrop() bool {
	return c.GetBool("cantChooseAfterDrop")
}

func (c *Cell) IsSafeDrop() bool {
	return c.GetBool("isSafeDrop")
}
