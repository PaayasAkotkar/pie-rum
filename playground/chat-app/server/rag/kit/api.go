package kit

import (
	"context"
	"sync"

	"github.com/firebase/genkit/go/genkit"
)

func (k *Kit[In, Out]) Add(g *genkit.Genkit, name string, run handler[In, Out]) *Kit[In, Out] {
	var once sync.Once
	once.Do(func() {
		f := genkit.DefineFlow(g, name, run)
		k.flows[name] = f
		k.keys.Push(name)
	})
	return k
}

func (k *Kit[In, Out]) Run(ctx context.Context, name string, in In) *Kit[In, Out] {
	f, err := k.flows[name].Run(ctx, in)
	o := &Output[Out]{
		Err: err,
		Res: f,
	}
	k.cheetah.Publish(name, o)
	return k
}

func (k *Kit[In, Out]) Monitor(ctx context.Context, name string, fn func(handler *Output[Out])) {
	ch := k.cheetah.Subscribe(name)
	go func() {
		select {
		case t := <-ch:
			fn(t)
		case <-ctx.Done():
			k.cheetah.Unsubscribe(name, ch)
			return
		}
	}()
}

func (k *Kit[In, Out]) AutoMonitor(ctx context.Context, fn func(handler *Output[Out])) {
	out := make(chan *Output[Out], 100)
	go func() {
		for {
			select {
			case result, ok := <-out:
				if !ok {
					return
				}
				if result == nil {
					continue
				}
				fn(result)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for _, r := range k.keys.Max() {
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					ch := k.cheetah.Subscribe(r)
					select {
					case <-ctx.Done():
						k.cheetah.Unsubscribe(r, ch)
						return
					case result, ok := <-ch:
						k.cheetah.Unsubscribe(r, ch)
						if !ok || result == nil {
							return
						}
						select {
						case out <- result:
						case <-ctx.Done():
							return
						}
					}
				}
			}()
		}
	}()
}
