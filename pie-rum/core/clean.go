package pierum

// clean releases the Stored metrics and all to free up space
func (r *PieRum[In, Out]) clean() {
	r.mu.Lock()
	defer r.mu.Unlock()
	go r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for _, p := range r.Store.Registry {
			for _, k := range p.Registry {
				for _, s := range k.Registry {
					for _, d := range s.Registry {
						d.release()
					}
				}
			}
		}
		r.Store.release()
	}()
	r.wg.Wait()
}
