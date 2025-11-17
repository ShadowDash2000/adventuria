package event

type UnsubGroup struct {
	Fns []Unsubscribe
}

func (g *UnsubGroup) Add(u ...Unsubscribe) {
	g.Fns = append(g.Fns, u...)
}

func (g *UnsubGroup) Unsubscribe() {
	for _, fn := range g.Fns {
		fn()
	}
	g.Fns = nil
}
