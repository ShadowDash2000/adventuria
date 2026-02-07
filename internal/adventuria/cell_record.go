package adventuria

import (
	"adventuria/internal/adventuria/schema"
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
	return c.GetInt(schema.CellSchema.Sort)
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
	return c.GetString(schema.CellSchema.Filter)
}

func (c *CellRecord) AudioPreset() string {
	return c.GetString(schema.CellSchema.AudioPreset)
}

func (c *CellRecord) Icon() string {
	return c.GetString(schema.CellSchema.Icon)
}

func (c *CellRecord) Name() string {
	return c.GetString(schema.CellSchema.Name)
}

func (c *CellRecord) Points() int {
	return c.GetInt(schema.CellSchema.Points)
}

func (c *CellRecord) Coins() int {
	return c.GetInt(schema.CellSchema.Coins)
}

func (c *CellRecord) Description() string {
	return c.GetString(schema.CellSchema.Description)
}

func (c *CellRecord) Color() string {
	return c.GetString(schema.CellSchema.Color)
}

func (c *CellRecord) CantDrop() bool {
	return c.GetBool(schema.CellSchema.CantDrop)
}

func (c *CellRecord) CantReroll() bool {
	return c.GetBool(schema.CellSchema.CantReroll)
}

func (c *CellRecord) IsSafeDrop() bool {
	return c.GetBool(schema.CellSchema.IsSafeDrop)
}

func (c *CellRecord) IsCustomFilterNotAllowed() bool {
	return c.GetBool(schema.CellSchema.IsCustomFilterNotAllowed)
}

func (c *CellRecord) Value() string {
	return c.GetString(schema.CellSchema.Value)
}
func (c *CellRecord) UnmarshalValue(result any) error {
	return c.UnmarshalJSONField(schema.CellSchema.Value, result)
}
