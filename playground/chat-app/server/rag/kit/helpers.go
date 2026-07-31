package kit

import "context"

type h[Out any] func() Out

type handler[In, Out any] func(ctx context.Context, in In) (Out, error)

type Output[O any] struct {
	Err     error
	IsReady bool
	Res     O
}
