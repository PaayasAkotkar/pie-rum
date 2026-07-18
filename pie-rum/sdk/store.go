package pierum

import (
	"context"
	rum "pie-rum-sdk/pie-rum/core"

	"sync"
)

type Store[In, Out any] struct {
	profile map[string]*Profile[In, Out]
	core    *rum.IStore[In, Out]
	mu      sync.Mutex
}

func NewStore[In, Out any](ctx context.Context) *Store[In, Out] {
	return &Store[In, Out]{
		profile: make(map[string]*Profile[In, Out]),
		core:    rum.NewStore[In, Out](ctx),
	}
}

func (s *Store[In, Out]) AddProfile(name string, p func(p *Profile[In, Out])) *Store[In, Out] {
	engine := NewNode(s, func() *Profile[In, Out] {
		return NewProfile[In, Out](len(s.profile) + 1)
	})
	return engine.Nest(name, s.profile, p)
}

func (s *Store[In, Out]) SetProfile(name string, rank int, p func(p *Profile[In, Out])) *Store[In, Out] {
	engine := NewNode(s, func() *Profile[In, Out] {
		return NewProfile[In, Out](rank)
	})

	return engine.Nest(name, s.profile, p)
}

func (s *Store[In, Out]) ReplaceProfile(name string, profile *Profile[In, Out]) *Store[In, Out] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profile[name] = profile
	s.core.ReplaceProfile(name, profile.core)
	return s
}

func (s *Store[In, Out]) GetCoreStore() *rum.IStore[In, Out] {
	return s.core
}

func (s *Store[In, Out]) Build() *Store[In, Out] {
	createTags := []string{}
	for name, p := range s.profile {
		p.Build() // build profile's children first
		s.core.AddProfile(name, p.core)
		createTags = append(createTags, name)
	}
	s.core.SetMonitorTags(createTags)

	return s
}
