package pierum

import (
	"log"
	"runtime/debug"
)

// Hub listens the channels request
func (r*PieRum[In, Out]) Hub() {
	log.Println("listening 🐴")

	for {
		func() {

			defer func() {
				if r := recover(); r != nil {
					log.Printf("PANIC RECOVERED in Hub: %v\n%s", r, debug.Stack())
				}
			}()

			select {

			case token := <-r.post:
				log.Println("in listening post")
				r.addRequestCount()
				for _, l := range token.Links {
					r.onPost(l.Seq)
				}

			case token := <-r.release:
				if token {
					r.clean()
					r.resetRequestCount()
				}

			case token := <-r.monitorTags:
				r.addRequestCount()
				r.store.SetMonitorTags(token)

			case tokens := <-r.activateProfile:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to activate profile: %v", err)
					}
				}

			case tokens := <-r.activateKit:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to activate kit: %v", err)
					}
				}

			case tokens := <-r.activateService:
				r.addRequestCount()

				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to activate service: %v", err)
					}
				}

			case tokens := <-r.activateDispatcher:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to activate dispatcher: %v", err)
					}
				}

			case tokens := <-r.activateEvent:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to activate event: %v", err)
					}
				}

			case tokens := <-r.deactivateProfile:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to deactivate profile: %v", err)
					}
				}

			case tokens := <-r.deactivateKit:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to deactivate kit: %v", err)
					}
				}

			case tokens := <-r.deactivateService:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to deactivate service: %v", err)
					}
				}

			case tokens := <-r.deactivateDispatcher:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to deactivate dispatcher: %v", err)
					}
				}

			case tokens := <-r.deactivateEvent:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to deactivate event: %v", err)
					}
				}

			case tokens := <-r.swapProfile:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to swap profile: %v", err)
					}
				}

			case tokens := <-r.swapKit:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to swap kit: %v", err)
					}
				}

			case tokens := <-r.swapService:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to swap service: %v", err)
					}
				}

			case tokens := <-r.swapDispatcher:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to swap dispatcher: %v", err)
					}
				}

			case tokens := <-r.swapEvent:
				r.addRequestCount()
				for _, token := range tokens {
					if err := r.handleActions(token); err != nil {
						log.Printf("failed to swap event: %v", err)
					}
				}

			case <-r.ctx.Done():
				// r.wg.Wait().
				printHeader()
				return
			}
		}()
	}
}
