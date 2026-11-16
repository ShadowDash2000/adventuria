package actors

import "time"

type ActiveActor struct {
	Id                string
	Kind              ActorKind
	Updated           time.Time
	Name              string
	Avatar            string
	Color             string
	CurrentWorld      string
	CellsPassed       int
	NextShiftChangeAt *time.Time
}
