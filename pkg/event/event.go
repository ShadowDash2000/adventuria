package event

type Resolver interface {
	Next() error

	setNextFunc(f func() error)
}

type Event struct {
	next func() error
}

func (e *Event) Next() error {
	if e.next != nil {
		return e.next()
	}
	return nil
}

func (e *Event) setNextFunc(fn func() error) {
	e.next = fn
}
