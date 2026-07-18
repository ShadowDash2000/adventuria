package model

import "context"

type CellType string

type Cell interface {
	Data() *CellInfo
	Categories() []string
	InCategory(category string) bool
	InCategories(categories []string) bool

	OnCellReached(ctx context.Context, events *Events, player *Player, reachedCtx *ReachedContext) error
	OnCellLeft(ctx context.Context, events *Events, player *Player) error
}

type ReachedContext struct {
	Moves []*MoveResult
}

type Verifiable interface {
	Verify(ctx context.Context, value string) error
}

type Rollable interface {
	Cell
	Refreshable
	Roll(ctx context.Context, events *Events, player *Player) (*WheelRollResult, error)
}

type Refreshable interface {
	RefreshItems(ctx context.Context, events *Events, player *Player) error
}

type CellData struct {
	Id                       string
	Disabled                 bool
	Sort                     int
	Type                     CellType
	World                    string
	Filter                   string
	AudioPreset              string
	Icon                     string
	Name                     string
	Points                   int
	EnergyConsume            int
	Coins                    int
	Description              string
	Color                    string
	CantDrop                 bool
	CantReroll               bool
	IsSafeDrop               bool
	IsCustomFilterNotAllowed bool
	IsChangeGameNotAllowed   bool
	DontGiveItemWheel        bool
	Value                    string

	LocalOrder  int
	GlobalOrder int
}

type CellInfo struct {
	data CellData
}

func RestoreCellInfo(data CellData) *CellInfo {
	return &CellInfo{data: data}
}

func (c *CellInfo) ID() string {
	return c.data.Id
}

func (c *CellInfo) Disabled() bool {
	return c.data.Disabled
}

func (c *CellInfo) Sort() int {
	return c.data.Sort
}

func (c *CellInfo) Type() CellType {
	return c.data.Type
}

func (c *CellInfo) World() string {
	return c.data.World
}

func (c *CellInfo) Filter() string {
	return c.data.Filter
}

func (c *CellInfo) AudioPreset() string {
	return c.data.AudioPreset
}

func (c *CellInfo) Icon() string {
	return c.data.Icon
}

func (c *CellInfo) Name() string {
	return c.data.Name
}

func (c *CellInfo) Points() int {
	return c.data.Points
}

func (c *CellInfo) EnergyConsume() int {
	return c.data.EnergyConsume
}

func (c *CellInfo) Coins() int {
	return c.data.Coins
}

func (c *CellInfo) Description() string {
	return c.data.Description
}

func (c *CellInfo) Color() string {
	return c.data.Color
}

func (c *CellInfo) CantDrop() bool {
	return c.data.CantDrop
}

func (c *CellInfo) CantReroll() bool {
	return c.data.CantReroll
}

func (c *CellInfo) IsSafeDrop() bool {
	return c.data.IsSafeDrop
}

func (c *CellInfo) IsCustomFilterNotAllowed() bool {
	return c.data.IsCustomFilterNotAllowed
}

func (c *CellInfo) IsChangeGameNotAllowed() bool {
	return c.data.IsChangeGameNotAllowed
}

func (c *CellInfo) DontGiveItemWheel() bool {
	return c.data.DontGiveItemWheel
}

func (c *CellInfo) Value() string {
	return c.data.Value
}

func (c *CellInfo) LocalOrder() int {
	return c.data.LocalOrder
}

func (c *CellInfo) GlobalOrder() int {
	return c.data.GlobalOrder
}
