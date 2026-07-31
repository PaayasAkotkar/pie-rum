package pierum

import (
 rum "pie-rum-sdk/pie-rum/core"
	"context"
)

type Event[In, Out any] struct {
	rank int
	core *rum.IEvent[In, Out]
}

func NewEvent[In, Out any](rank int) *Event[In, Out] {
	return &Event[In, Out]{
		rank: rank,
		core: rum.NewRegisterFunc[In, Out](),
	}
}

func (e *Event[In, Out]) RegisterFunc(fn func(ctx context.Context, in In) (Out, error)) *Event[In, Out] {
	e.core.Fn = fn
	e.core.Rank = int64(e.rank)
	return e
}
