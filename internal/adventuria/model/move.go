package model

type MoveResult struct {
	Type           MoveType  `json:"type"`
	Steps          int       `json:"steps"`
	TotalSteps     int       `json:"total_steps"`
	PrevTotalSteps int       `json:"prev_total_steps"`
	CurrentCell    *CellInfo `json:"current_cell"`
	CurrentWorld   *World    `json:"current_world"`
	Laps           int       `json:"laps"`
}

type MoveType string

const (
	MoveTypePath            MoveType = "path"
	MoveTypeTeleport        MoveType = "teleport"
	MoveTypeWorldTransition MoveType = "world_transition"
)
