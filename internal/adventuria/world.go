package adventuria

import "github.com/pocketbase/pocketbase/core"

type World interface {
	core.RecordProxy

	ID() string
	Name() string
	Slug() string
	Sort() int
	IsLoop() bool
	IsDefaultWorld() bool
	TransitionToWorld() string
	Effects() []string
}
