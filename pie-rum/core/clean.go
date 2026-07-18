package pierum

// clean releases the stored metrics and all to free up space
func (r*PieRum[In, Out]) clean() {
	r.mu.Lock()
	defer r.mu.Unlock()
	go r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for _, p := range r.store.registry {
			for _, k := range p.registry {
				for _, s := range k.registry {
					for _, d := range s.registry {
						d.release()
					}
				}
			}
		}
		r.store.release()
	}()
	r.wg.Wait()
}
