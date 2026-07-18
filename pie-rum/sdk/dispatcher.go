package pierum

import rum "pie-rum-sdk/pie-rum/core"

type Dispatcher[In, Out any] struct {
	rank   int
	name   string
	events map[string]*Event[In, Out]
	core   *rum.IDispatcher[In, Out]
}

func NewDispatcher[In, Out any](rank int) *Dispatcher[In, Out] {
	return &Dispatcher[In, Out]{
		rank:   rank,
		events: make(map[string]*Event[In, Out]),
		core:   rum.NewDispatcher[In, Out](rum.Settings{}), // name is set in build
	}
}

func (d *Dispatcher[In, Out]) AddEvent(name string, fn func(event *Event[In, Out])) *Dispatcher[In, Out] {
	engine := NewNode(d, func() *Event[In, Out] {
		return NewEvent[In, Out](len(d.events) + 1)
	})
	return engine.Nest(name, d.events, fn)
}

func (d *Dispatcher[In, Out]) ReplaceEvent(name string, event *Event[In, Out]) *Dispatcher[In, Out] {
	d.core.ReplaceEvent(name, event.core)
	return d
}

func (d *Dispatcher[In, Out]) SetEvent(name string, rank int, fn func(event *Event[In, Out])) *Dispatcher[In, Out] {
	engine := NewNode(d, func() *Event[In, Out] {
		return NewEvent[In, Out](rank)
	})
	return engine.Nest(name, d.events, fn)
}

func (d *Dispatcher[In, Out]) Build() *Dispatcher[In, Out] {
	for n, r := range d.events {
		r.core.Rank = int64(r.rank)
		d.core.PushEvent(n, r.core)
		delete(d.events, n)
	}
	return d
}
