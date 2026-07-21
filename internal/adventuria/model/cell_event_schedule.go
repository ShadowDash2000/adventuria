package model

import "time"

type CellEventScheduleData struct {
	Id              string
	ActionEvent     string
	Effects         []string
	ActiveCell      string
	CellTypes       []CellType
	Worlds          []string
	ShiftInterval   int
	LastShiftChange time.Time
}

type CellEventSchedule struct {
	data CellEventScheduleData
}

func RestoreCellEventSchedule(data CellEventScheduleData) *CellEventSchedule {
	return &CellEventSchedule{data: data}
}

func (c *CellEventSchedule) ID() string {
	return c.data.Id
}

func (c *CellEventSchedule) ActionEvent() string {
	return c.data.ActionEvent
}

func (c *CellEventSchedule) Effects() []string {
	return c.data.Effects
}

func (c *CellEventSchedule) ActiveCell() string {
	return c.data.ActiveCell
}

func (c *CellEventSchedule) SetActiveCell(cell string) {
	c.data.ActiveCell = cell
}

func (c *CellEventSchedule) CellTypes() []CellType {
	return c.data.CellTypes
}

func (c *CellEventSchedule) Worlds() []string {
	return c.data.Worlds
}

func (c *CellEventSchedule) ShiftInterval() int {
	return c.data.ShiftInterval
}

func (c *CellEventSchedule) LastShiftChange() time.Time {
	return c.data.LastShiftChange
}

func (c *CellEventSchedule) SetLastShiftChange(t time.Time) {
	c.data.LastShiftChange = t
}
