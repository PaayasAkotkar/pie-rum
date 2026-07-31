package injection

import "log"

func (h *Injection) pipe() {
	defer close(h.done)

	container := NewContainer(h.ctx, h.nodeID)
	h.container = container

	for {
		select {
		case req := <-h.getService:
			h.handleGetService(req)

		case reg := <-h.addService:
			h.handleAddService(reg, container)

		case resultCh := <-h.buildServices:
			resultCh <- h.handleBuild(container)

		case <-h.rebuildSignal:
			log.Printf("[%s] Rebuild signal received, reinitializing DI\n", h.nodeID)
			container = NewContainer(h.ctx, h.nodeID)
			h.container = container

		case <-h.stopChan:
			return
		case <-h.ctx.Done():
			return
		}
	}
}
