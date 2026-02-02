package adventuria

import (
	"slices"

	"github.com/pocketbase/pocketbase/core"
)

type CellRecord struct {
	core.BaseRecordProxy
	t CellType
}

func (c *CellRecord) ID() string {
	return c.Id
}

func (c *CellRecord) Sort() int {
	return c.GetInt("sort")
}

func (c *CellRecord) Type() CellType {
	return c.t
}

func (c *CellRecord) setType(t CellType) {
	c.t = t
}

func (c *CellRecord) Categories() []string {
	if def, ok := cellsList[c.Type()]; ok {
		return def.Categories
	}

	return nil
}

func (c *CellRecord) InCategory(category string) bool {
	return slices.Contains(c.Categories(), category)
}

func (c *CellRecord) Filter() string {
	return c.GetString("filter")
}

func (c *CellRecord) AudioPresets() []string {
	return c.GetStringSlice("audioPresets")
}

func (c *CellRecord) Icon() string {
	return c.GetString("icon")
}

func (c *CellRecord) Name() string {
	return c.GetString("name")
}

func (c *CellRecord) Points() int {
	return c.GetInt("points")
}

func (c *CellRecord) Coins() int {
	return c.GetInt("coins")
}

func (c *CellRecord) Description() string {
	return c.GetString("description")
}

func (c *CellRecord) Color() string {
	return c.GetString("color")
}

func (c *CellRecord) CantDrop() bool {
	return c.GetBool("cantDrop")
}

func (c *CellRecord) CantReroll() bool {
	return c.GetBool("cantReroll")
}

func (c *CellRecord) IsSafeDrop() bool {
	return c.GetBool("isSafeDrop")
}

func (c *CellRecord) IsCustomFilterNotAllowed() bool {
	return c.GetBool("is_custom_filter_not_allowed")
}

func (c *CellRecord) Value() string {
	return c.GetString("value")
}
