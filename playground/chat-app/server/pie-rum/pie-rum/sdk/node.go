package pierum

// Node implements the parent(P), child(C) mapping
type Node[P, C any] struct {
	par *P
	alc func() *C
}

func NewNode[P, C any](parent *P, allocate func() *C) *Node[P, C] {
	return &Node[P, C]{
		par: parent,
		alc: allocate,
	}
}

func (n *Node[P, C]) Nest(name string, registry map[string]*C, scope func(child *C)) *P {
	child := n.alc()
	scope(child)
	registry[name] = child
	return n.par
}
