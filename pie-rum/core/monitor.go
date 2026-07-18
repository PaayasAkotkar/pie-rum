package pierum

func (r *PieRum[In, Out]) fetch(profile string) *IResults {
	ch := r.cheetah.Subscribe(profile)
	defer r.cheetah.Unsubscribe(profile, ch)

	select {
	case <-r.ctx.Done():
		return nil
	case result := <-ch:
		return result
	}
}
