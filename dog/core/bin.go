package dog

func (rd *Dog[T]) bin(req IDone) {
	p, exists := rd.policy[req.PolicyName]
	if !exists {
		return
	}

	newFns := make([]Funcs[T], 0, len(p.Fn))
	for _, r := range p.Fn {
		if r.Rank != req.Rank {
			newFns = append(newFns, r)
		}
	}
	p.Fn = newFns
}
