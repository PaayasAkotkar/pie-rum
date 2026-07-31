package kit

import (
	"pie-rum-sdk/cheetah"
	"pie-rum-sdk/stack"

	"github.com/firebase/genkit/go/core"
)

type Kit[In, Out any] struct {
	flows   map[string]*core.Flow[In, Out, struct{}]
	keys    *stack.Stack[string]
	cheetah *cheetah.Cheetah[string, Output[Out]]
}

func New[In, Out any]() *Kit[In, Out] {
	return &Kit[In, Out]{
		flows:   make(map[string]*core.Flow[In, Out, struct{}]),
		cheetah: cheetah.New[string, Output[Out]](100),
		keys:    stack.NewStack[string](),
	}
}
