// package cheetah

// import (
// 	"sync"
// )

// // Cheetah  is a generic pub/sub system keyed by string.
// // T is the message payload type.
// type Cheetah[T any] struct {
// 	mu          sync.Mutex
// 	subscribers map[string]map[chan *T]struct{}
// }

// // New returns new Cheetah  of single buffer
// // single for fast process 😄
// func New[T any](buffer int) *Cheetah[T] {
// 	return &Cheetah[T]{
// 		subscribers: make(map[string]map[chan *T]struct{}, buffer),
// 	}
// }

// // Subscribe returns a buffered channel that receives published values for key
// func (l *Cheetah[T]) Subscribe(key string) chan *T {
// 	ch := make(chan *T, 1)
// 	l.mu.Lock()
// 	defer l.mu.Unlock()
// 	if _, ok := l.subscribers[key]; !ok {
// 		l.subscribers[key] = make(map[chan *T]struct{}, 1)
// 	}
// 	l.subscribers[key][ch] = struct{}{}
// 	return ch
// }

// // Publish sends results to all subscribers of key. Never blocks.
// func (l *Cheetah[T]) Publish(key string, parcel *T) {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()
// 	for ch := range l.subscribers[key] {
// 		select {
// 		case ch <- parcel:
// 		default:
// 		}
// 	}
// }

// // Unsubscribe removes and closes the channel for key
//
//	func (l *Cheetah[T]) Unsubscribe(key string, body chan *T) {
//		l.mu.Lock()
//		defer l.mu.Unlock()
//		if sub, ok := l.subscribers[key]; ok {
//			delete(sub, body)
//			close(body)
//			if len(sub) == 0 {
//				delete(l.subscribers, key)
//			}
//		}
//	}
package cheetah

import (
	"sync"
)

// Cheetah is a generic pub/sub system keyed by string.
// T is the message payload type.
type Cheetah[K comparable, T any] struct {
	mu          sync.Mutex
	subscribers map[K]map[chan *T]struct{}
	last        map[K]*T // ← latched last value per key
}

// New returns a new Cheetah instance.
func New[K comparable, T any](buffer int) *Cheetah[K, T] {
	return &Cheetah[K, T]{
		subscribers: make(map[K]map[chan *T]struct{}, buffer),
		last:        make(map[K]*T),
	}
}

// Subscribe returns a buffered channel that receives published values for key.
// If a value was already published for this key, it is replayed immediately.
// func (l *Cheetah[K, T]) Subscribe(key K) chan *T {
// 	ch := make(chan *T, 1)
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	// Replay last value to late subscriber — solves the race
// 	if val, ok := l.last[key]; ok {
// 		ch <- val
// 		return ch // no need to register, already has its answer
// 	}

//		if _, ok := l.subscribers[key]; !ok {
//			l.subscribers[key] = make(map[chan *T]struct{}, 1)
//		}
//		l.subscribers[key][ch] = struct{}{}
//		return ch
//	}
func (l *Cheetah[K, T]) Subscribe(key K) chan *T {
	ch := make(chan *T, 1)
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.subscribers[key]; !ok {
		l.subscribers[key] = make(map[chan *T]struct{}, 1)
	}
	l.subscribers[key][ch] = struct{}{}

	if val, ok := l.last[key]; ok {
		ch <- val
		delete(l.last, key) // ← evict here, only one subscriber gets the replay
	}

	return ch
}

// Publish sends results to all subscribers of key and latches the value
// for any future late subscribers. Never blocks.
// func (l *Cheetah[K, T]) Publish(key K, parcel *T) {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	l.last[key] = parcel // latch before notifying

//		for ch := range l.subscribers[key] {
//			select {
//			case ch <- parcel:
//			default:
//			}
//		}
//	}

func (l *Cheetah[K, T]) Publish(key K, parcel *T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.last[key] = parcel

	for ch := range l.subscribers[key] {
		select {
		case ch <- parcel:
		default:
		}
	}
	// if this is not done the polling keeps getting the request
	delete(l.subscribers, key)
	// delete(l.last, key)
}

// Unsubscribe removes and closes the channel for key.
func (l *Cheetah[K, T]) Unsubscribe(key K, body chan *T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if sub, ok := l.subscribers[key]; ok {
		delete(sub, body)
		close(body)
		if len(sub) == 0 {
			delete(l.subscribers, key)
		}
	}
}

// Evict clears the latched value for key (call after result is consumed).
func (l *Cheetah[K, T]) Evict(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}
