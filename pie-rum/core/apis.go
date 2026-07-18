// Package pierum the code is improved alot by google's gemini so shotout to gemini
package pierum

import (
	"context"
	"log"
)

// Paper — one-shot, closes after first value, never fetches again
func (r *PieRum[In, Out]) Paper(profile string) <-chan *IResults {
	out := make(chan *IResults, 1)
	go func() {
		defer close(out)
		result := r.fetch(profile)
		out <- result
	}()
	return out
}

// GoPoll — starts server fns + streams results to handler
// based on the request key you can specify which service to run
// handler is called once per publish, parks between events
func (r *PieRum[In, Out]) GoPoll(ctx context.Context, profile string, handler Handler, fns ...func()) {
	pollCh := r.Poll(profile)

	for _, fn := range fns {
		fn := fn
		go fn()
	}

	go func() {
		for result := range pollCh {
			if result == nil {
				continue
			}
			handler(result)
		}
	}()
}

// GoAuto — starts server fns + streams results to handler
// will only monitor all the registered profiles
func (r *PieRum[In, Out]) GoAuto(ctx context.Context, handler Handler, fns ...func()) {
	out := make(chan *IResults, len(r.store.registry))

	// fan-in all registered profiles into one channel
	for profile := range r.store.registry {
		profile := profile
		go func() {
			ch := r.Poll(profile)
			for result := range ch {
				select {
				case <-r.ctx.Done():
					return
				case out <- result:
				}
			}
		}()
	}

	for _, fn := range fns {
		fn := fn
		go fn()
	}

	go func() {
		defer close(out)
		<-r.ctx.Done()
	}()

	go func() {
		for result := range out {
			if result == nil {
				continue
			}
			handler(result)
		}
	}()
}

// Poll — sse-style, one event delivered then parks, lives until ctx cancel
func (r *PieRum[In, Out]) Poll(profile string) <-chan *IResults {
	out := make(chan *IResults, 1)
	go func() {
		defer close(out)
		for {
			ch := r.cheetah.Subscribe(profile)
			select {
			case <-r.ctx.Done():
				r.cheetah.Unsubscribe(profile, ch)
				return
			case result, ok := <-ch:
				r.cheetah.Unsubscribe(profile, ch)
				if !ok || result == nil {
					return
				}
				select {
				case out <- result:
				case <-r.ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// GoMonitor listens to the results of the hub methods
func (r *PieRum[In, Out]) GoMonitor(ctx context.Context, handler Handler, fns ...func()) {
	out := make(chan *IResults, buffers)

	go func() {
		for {
			select {
			case result, ok := <-out:
				if !ok {
					return
				}
				if result == nil {
					continue
				}
				log.Println("results: ", result.Resuts)
				handler(result)
			case <-r.ctx.Done():
				return
			}
		}
	}()

	for _, profile := range r.store.GetTags() {
		profile := profile
		go func() {
			for {
				select {
				case <-r.ctx.Done():
					return
				default:
				}

				ch := r.cheetah.Subscribe(profile)
				select {
				case <-r.ctx.Done():
					r.cheetah.Unsubscribe(profile, ch)
					return
				case result, ok := <-ch:
					r.cheetah.Unsubscribe(profile, ch)
					if !ok || result == nil {
						return
					}
					select {
					case out <- result:
					case <-r.ctx.Done():
						return
					}
				}
			}
		}()
	}

	for _, fn := range fns {
		fn := fn
		go fn()
	}
}

// GoErrors — keyless, receives ALL errors published by any autoWrite call.
// handler receives the raw error; no key knowledge required.
func (r *PieRum[In, Out]) GoErrors(ctx context.Context, handler func(err error), fns ...func()) {
	out := make(chan error, 8)

	go func() {
		ch := r.cheetahDetector.Subscribe("*")
		defer r.cheetahDetector.Unsubscribe("*", ch)
		for {
			select {
			case <-r.ctx.Done():
				return
			case e, ok := <-ch:
				if !ok || e == nil {
					return
				}
				select {
				case out <- *e:
				case <-r.ctx.Done():
					return
				}
			}
		}
	}()

	for _, fn := range fns {
		fn := fn
		go fn()
	}

	go func() {
		for err := range out {
			handler(err)
		}
	}()
}

// GoErrorsKeyed — for callers who DO know the key and want targeted listening.
// Key format: "profile.kit.service.dispatcher.event"
func (r *PieRum[In, Out]) GoErrorsKeyed(ctx context.Context, key string, handler func(err error), fns ...func()) {
	out := make(chan error, 8)

	go func() {
		ch := r.cheetahDetector.Subscribe(key)
		defer r.cheetahDetector.Unsubscribe(key, ch)
		for {
			select {
			case <-r.ctx.Done():
				return
			case e, ok := <-ch:
				if !ok || e == nil {
					return
				}
				select {
				case out <- *e:
				case <-r.ctx.Done():
					return
				}
			}
		}
	}()

	for _, fn := range fns {
		fn := fn
		go fn()
	}

	go func() {
		for err := range out {
			handler(err)
		}
	}()
}
func (r *PieRum[In, Out]) GetStore() *IStore[In, Out] {
	return r.store
}
func (r *PieRum[In, Out]) GetProfiles() map[string]*IProfile[In, Out] {
	return r.store.registry
}
func (r *PieRum[In, Out]) GetKits(profile string) map[string]*IKit[In, Out] {
	return r.store.registry[profile].registry
}
func (r *PieRum[In, Out]) GetServices(profile, kit string) map[string]*IService[In, Out] {
	return r.store.registry[profile].registry[kit].registry
}
func (r *PieRum[In, Out]) GetDispatchers(profile, kit, service string) map[string]*IDispatcher[In, Out] {
	return r.store.registry[profile].registry[kit].registry[service].registry
}
func (r *PieRum[In, Out]) GetEvents(profile, kit, service, dispatcher string) map[string]*IEvent[In, Out] {
	return r.store.registry[profile].registry[kit].registry[service].registry[dispatcher].registry
}
