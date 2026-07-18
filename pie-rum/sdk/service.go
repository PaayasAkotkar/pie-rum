package pierum

import (
	 rum "pie-rum-sdk/pie-rum/core"

	"context"
)

type Service[In, Out any] struct {
	name        string
	rank        int
	core        *rum.IService[In, Out]
	dispatchers map[string]*Dispatcher[In, Out]
}

func NewService[In, Out any](rank int) *Service[In, Out] {
	return &Service[In, Out]{
		rank:        rank,
		dispatchers: make(map[string]*Dispatcher[In, Out]),
		core:        rum.NewService[In, Out](context.Background(), rum.Settings{}, ""), // Name is set in Build
	}
}

func (s *Service[In, Out]) AddDispatcher(name string, d func(d *Dispatcher[In, Out])) *Service[In, Out] {
	engine := NewNode(s, func() *Dispatcher[In, Out] {
		return NewDispatcher[In, Out](len(s.dispatchers) + 1)
	})
	return engine.Nest(name, s.dispatchers, d)
}
func (s *Service[In, Out]) SetDispatcher(name string, rank int, fn func(event *Dispatcher[In, Out])) *Service[In, Out] {
	engine := NewNode(s, func() *Dispatcher[In, Out] {
		return NewDispatcher[In, Out](rank)
	})
	return engine.Nest(name, s.dispatchers, fn)
}

func (s *Service[In, Out]) Build() *Service[In, Out] {
	for n, r := range s.dispatchers {
		r.Build() 
		r.core.Name = n
		r.core.Rank = int64(r.rank)
		s.core.PushDispatcher(n, r.core)
		delete(s.dispatchers, n)
	}
	return s
}

func (s *Service[In, Out]) ReplaceDispatcher(name string, dp *Dispatcher[In, Out]) *Service[In, Out] {
	s.core.ReplaceDispatcher(name, dp.core)
	return s
}
