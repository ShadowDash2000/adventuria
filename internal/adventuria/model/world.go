package model

type WorldData struct {
	Id                string
	Name              string
	Slug              string
	Sort              int
	IsLoop            bool
	IsDefaultWorld    bool
	TransitionToWorld string
	Effects           []string
}

type World struct {
	data  WorldData
	isNew bool
}

func RestoreWorld(data WorldData) *World {
	return &World{
		data:  data,
		isNew: false,
	}
}

func (w *World) IsNew() bool {
	return w.isNew
}

func (w *World) ID() string {
	return w.data.Id
}

func (w *World) Name() string {
	return w.data.Name
}

func (w *World) Slug() string {
	return w.data.Slug
}

func (w *World) Sort() int {
	return w.data.Sort
}

func (w *World) IsLoop() bool {
	return w.data.IsLoop
}

func (w *World) IsDefaultWorld() bool {
	return w.data.IsDefaultWorld
}

func (w *World) TransitionToWorld() string {
	return w.data.TransitionToWorld
}

func (w *World) Effects() []string {
	return w.data.Effects
}
