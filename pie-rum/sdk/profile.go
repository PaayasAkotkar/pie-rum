package pierum

import rum "pie-rum-sdk/pie-rum/core"

type Profile[In, Out any] struct {
	rank int
	name string
	kits map[string]*Kit[In, Out] // temp will be nil after build
	core *rum.IProfile[In, Out]
}

func NewProfile[In, Out any](rank int) *Profile[In, Out] {
	return &Profile[In, Out]{
		rank: rank,
		kits: make(map[string]*Kit[In, Out]),
		core: rum.NewProfile[In, Out](),
	}
}

func (p *Profile[In, Out]) AddKit(name string, kit func(kit *Kit[In, Out])) *Profile[In, Out] {
	engine := NewNode(p, func() *Kit[In, Out] {
		return NewKit[In, Out](len(p.kits) + 1)
	})
	return engine.Nest(name, p.kits, kit)
}

func (p *Profile[In, Out]) ReplaceKit(name string, kit *Kit[In, Out]) *Profile[In, Out] {
	p.core.ReplaceKit(name, kit.core)
	return p
}

func (p *Profile[In, Out]) SetKit(name string, rank int, fn func(kit *Kit[In, Out])) *Profile[In, Out] {
	engine := NewNode(p, func() *Kit[In, Out] {
		return NewKit[In, Out](rank)
	})
	return engine.Nest(name, p.kits, fn)
}

func (p *Profile[In, Out]) Build() *Profile[In, Out] {
	for n, r := range p.kits {
		r.Build()
		r.core.Name = n
		r.core.Rank = int64(r.rank)
		p.core.PushKit(n, r.core)
		delete(p.kits, n)
	}
	return p
}
