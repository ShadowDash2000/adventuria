package event

import (
	"adventuria/pkg/random"
	"adventuria/pkg/result"
	"slices"
)

type Handler[T Resolver] struct {
	Func     func(T) (*result.Result, error)
	id       string
	once     bool
	Priority int
}

type Hook[T Resolver] struct {
	handlers []*Handler[T]
}

type Unsubscribe func()

func (h *Hook[T]) Bind(handler *Handler[T]) Unsubscribe {
	handler.id = generateHookId()

	insertAt := len(h.handlers)
	for i, existing := range h.handlers {
		if handler.Priority < existing.Priority {
			insertAt = i
			break
		}
	}

	h.handlers = slices.Insert(h.handlers, insertAt, handler)

	return func() {
		h.Unbind(handler.id)
	}
}

func (h *Hook[T]) BindFunc(fn func(e T) (*result.Result, error)) Unsubscribe {
	return h.Bind(&Handler[T]{
		Func: fn,
	})
}

func (h *Hook[T]) BindFuncWithPriority(fn func(e T) (*result.Result, error), priority int) Unsubscribe {
	return h.Bind(&Handler[T]{
		Func:     fn,
		Priority: priority,
	})
}

func (h *Hook[T]) BindFuncOnce(fn func(e T) (*result.Result, error)) Unsubscribe {
	return h.Bind(&Handler[T]{
		Func: fn,
		once: true,
	})
}

func (h *Hook[T]) Unbind(idsToRemove ...string) {
	for _, id := range idsToRemove {
		for i := len(h.handlers) - 1; i >= 0; i-- {
			if h.handlers[i].id == id {
				h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
				break
			}
		}
	}
}

func (h *Hook[T]) Trigger(event T, oneOffHandlerFuncs ...func(T) (*result.Result, error)) (*result.Result, error) {
	handlers := make([]*Handler[T], 0, len(h.handlers))
	handlers = append(handlers, h.handlers...)
	for _, fn := range oneOffHandlerFuncs {
		handlers = append(handlers, &Handler[T]{Func: fn})
	}

	event.setNextFunc(nil) // reset in case the event is being reused

	var onceIds []string

	for i := len(handlers) - 1; i >= 0; i-- {
		handler := handlers[i]
		old := event.nextFunc()
		event.setNextFunc(func() (*result.Result, error) {
			if handler.once && handler.id != "" {
				onceIds = append(onceIds, handler.id)
			}
			event.setNextFunc(old)
			return handler.Func(event)
		})

	}

	res, err := event.Next()

	if len(onceIds) > 0 {
		h.Unbind(onceIds...)
	}

	return res, err
}

func generateHookId() string {
	return random.PseudorandomString(20)
}
