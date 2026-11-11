package event

import (
	"context"

	"ganigeorgiev.com/pocketbase/pocketbase/tools/list"
)

type Tagger interface {
	Resolver

	Tags() []string
}

type mainHook[T Tagger] struct {
	*Hook[T]
}

func NewTaggedHook[T Tagger](hook *Hook[T], tags ...string) *TaggedHook[T] {
	return &TaggedHook[T]{
		mainHook[T]{hook},
		tags,
	}
}

type TaggedHook[T Tagger] struct {
	mainHook[T]

	tags []string
}

func (h *TaggedHook[T]) CanTriggerOn(tagsToCheck []string) bool {
	if len(h.tags) == 0 {
		return true // match all
	}

	for _, t := range tagsToCheck {
		if list.ExistInSlice(t, h.tags) {
			return true
		}
	}

	return false
}

func (h *TaggedHook[T]) Bind(handler *Handler[T]) Unsubscribe {
	fn := handler.Func

	handler.Func = func(ctx context.Context, e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(ctx, e)
		}

		return e.Next()
	}

	return h.mainHook.Bind(handler)
}

func (h *TaggedHook[T]) BindFunc(fn func(ctx context.Context, e T) error) Unsubscribe {
	return h.mainHook.BindFunc(func(ctx context.Context, e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(ctx, e)
		}

		return e.Next()
	})
}

func (h *TaggedHook[T]) BindFuncWithPriority(fn func(ctx context.Context, e T) error, priority int) Unsubscribe {
	return h.mainHook.BindFuncWithPriority(func(ctx context.Context, e T) error {
		if h.CanTriggerOn(e.Tags()) {
			return fn(ctx, e)
		}

		return e.Next()
	}, priority)
}
