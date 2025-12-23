package event

type Resolver interface {
	Next() (*Result, error)

	nextFunc() func() (*Result, error)
	setNextFunc(f func() (*Result, error))
}

type Result struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Event struct {
	next func() (*Result, error)
}

// Next calls the next hook handler.
func (e *Event) Next() (*Result, error) {
	if e.next != nil {
		return e.next()
	}
	return &Result{
		Success: true,
	}, nil
}

// nextFunc returns the function that Next calls.
func (e *Event) nextFunc() func() (*Result, error) {
	return e.next
}

// setNextFunc sets the function that Next calls.
func (e *Event) setNextFunc(fn func() (*Result, error)) {
	e.next = fn
}
