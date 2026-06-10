package model

type MoveResult struct {
	Steps          int       `json:"steps"`
	TotalSteps     int       `json:"total_steps"`
	PrevTotalSteps int       `json:"prev_total_steps"`
	CurrentCell    *CellInfo `json:"current_cell"`
	CellLocalOrder int       `json:"cell_local_order"`
	CurrentWorld   *World    `json:"current_world"`
	Laps           int       `json:"laps"`
}
