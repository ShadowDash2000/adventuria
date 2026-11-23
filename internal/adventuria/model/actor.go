package model

type ActorPosition struct {
	WorldId        string
	CellId         string
	CellLocalOrder int
}

type ActorMovement struct {
	Position ActorPosition
	Steps    int
	Laps     int
}
