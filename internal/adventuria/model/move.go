package model

type MoveResult struct {
	Type           MoveType
	TotalSteps     int
	PrevTotalSteps int
	World          *World
	Cell           *CellInfo
	Laps           int
}

type MoveType string

const (
	MoveTypePath            MoveType = "path"
	MoveTypeTeleport        MoveType = "teleport"
	MoveTypeWorldTransition MoveType = "world_transition"
)
