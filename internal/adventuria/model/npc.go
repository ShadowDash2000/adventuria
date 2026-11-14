package model

import "time"

type NPCData struct {
	Id                string
	Behavior          string
	ActionEvent       string
	Effects           []string
	CurrentWorld      string
	CurrentCell       string
	CellsPassed       int
	ShiftInterval     int
	NextShiftChangeAt time.Time
}

type NPC struct {
	data NPCData
}

func RestoreNPC(data NPCData) *NPC {
	return &NPC{data: data}
}

func (c *NPC) ID() string {
	return c.data.Id
}

func (c *NPC) Behavior() string {
	return c.data.Behavior
}

func (c *NPC) ActionEvent() string {
	return c.data.ActionEvent
}

func (c *NPC) Effects() []string {
	return c.data.Effects
}

func (c *NPC) CurrentWorld() string {
	return c.data.CurrentWorld
}

func (c *NPC) SetCurrentWorld(world string) {
	c.data.CurrentWorld = world
}

func (c *NPC) CurrentCell() string {
	return c.data.CurrentCell
}

func (c *NPC) SetCurrentCell(cell string) {
	c.data.CurrentCell = cell
}

func (c *NPC) CellsPassed() int {
	return c.data.CellsPassed
}

func (c *NPC) SetCellsPassed(n int) {
	c.data.CellsPassed = n
}

func (c *NPC) ShiftInterval() int {
	return c.data.ShiftInterval
}

func (c *NPC) NextShiftChangeAt() time.Time {
	return c.data.NextShiftChangeAt
}

func (c *NPC) SetNextShiftChangeAt(t time.Time) {
	c.data.NextShiftChangeAt = t
}
