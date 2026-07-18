package pierum

import rum "pie-rum-sdk/pie-rum/core"

type Kit[In, Out any] struct {
	rank     int
	name     string
	services map[string]*Service[In, Out]
	core     *rum.IKit[In, Out]
}

func NewKit[In, Out any](rank int) *Kit[In, Out] {
	return &Kit[In, Out]{
		rank:     rank,
		services: make(map[string]*Service[In, Out]),
		core:     rum.NewKit[In, Out](),
	}
}

func (k *Kit[In, Out]) AddService(name string, service func(service *Service[In, Out])) *Kit[In, Out] {
	engine := NewNode(k, func() *Service[In, Out] {
		return NewService[In, Out](len(k.services) + 1)
	})
	return engine.Nest(name, k.services, service)
}
func (k *Kit[In, Out]) ReplaceService(name string, service *Service[In, Out]) *Kit[In, Out] {
	k.core.ReplaceService(name, service.core)
	return k
}
func (k *Kit[In, Out]) SetService(name string, rank int, fn func(service *Service[In, Out])) *Kit[In, Out] {
	engine := NewNode(k, func() *Service[In, Out] {
		return NewService[In, Out](rank)
	})
	return engine.Nest(name, k.services, fn)
}
func (k *Kit[In, Out]) Build() *Kit[In, Out] {
	for n, r := range k.services {
		r.Build()
		r.core.Name = n
		r.core.Rank = int64(r.rank)
		k.core.PushService(n, r.core)
		delete(k.services, n)
	}
	return k
}
