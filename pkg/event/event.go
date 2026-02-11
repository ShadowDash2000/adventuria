package event

import "adventuria/pkg/result"

type Resolver interface {
	Next() (*result.Result, error)

	nextFunc() func() (*result.Result, error)
	setNextFunc(f func() (*result.Result, error))
}

type Event struct {
	next func() (*result.Result, error)
}

// Next calls the next hook handler.
func (e *Event) Next() (*result.Result, error) {
	if e.next != nil {
		return e.next()
	}
	return result.Ok(), nil
}

// nextFunc returns the function that Next calls.
func (e *Event) nextFunc() func() (*result.Result, error) {
	return e.next
}

// setNextFunc sets the function that Next calls.
func (e *Event) setNextFunc(fn func() (*result.Result, error)) {
	e.next = fn
}
