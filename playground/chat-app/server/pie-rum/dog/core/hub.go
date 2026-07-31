package dog

import (
	"log"
)

// watchDog is the main event loop for the Dog watchdog
func (rd *Dog[T]) watchDog() {
	defer rd.wg.Done()

	log.Println("🐕 Watchdog started...")

	for {
		select {
		// Registration handling
		case policy := <-rd.register:
			rd.handleRegister(policy)

		// Unregistration handling
		case name := <-rd.unregister:
			rd.handleUnregister(name)

		// Park/Monitor handling
		case policyName := <-rd.parkDog:
			rd.handleParkDog(policyName)

		// Summon/Execute handling
		case policyName := <-rd.summonCh:
			rd.handleSummon(policyName)

		// Completion handling
		case done := <-rd.done:
			rd.handleDone(done)

		// Error handling
		case bark := <-rd.bark:
			rd.handleBark(bark)

		// Reset handling
		case policyName := <-rd.reset:
			rd.handleReset(policyName)

		// Reset all handling
		case <-rd.resetAll:
			rd.handleResetAll()

		// Shutdown
		case <-rd.stopCh:
			log.Println("🛑 Watchdog shutting down...")
			return

		case <-rd.ctx.Done():
			log.Println("🛑 Context cancelled, shutting down...")
			return
		}
	}
}
