package dog

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PolicyState represents the lifecycle state of a policy
type PolicyState string

const (
	StateUnregistered PolicyState = "unregistered"
	StateRegistered   PolicyState = "registered"
	StateMonitoring   PolicyState = "monitoring"
	StateCompleted    PolicyState = "completed"
	StateError        PolicyState = "error"
)

// PolicyLifecycle tracks the state of a policy
type PolicyLifecycle struct {
	Name      string
	State     PolicyState
	mu        sync.RWMutex
	Timestamp time.Time
}

func (pl *PolicyLifecycle) GetState() PolicyState {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	return pl.State
}

func (pl *PolicyLifecycle) SetState(state PolicyState) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.State = state
	pl.Timestamp = time.Now()
}

func (pl *PolicyLifecycle) IsInState(state PolicyState) bool {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	return pl.State == state
}

// MonitorPolicy monitors execution of a policy
type MonitorPolicy struct {
	policy    string
	fn        func()
	stopChan  chan struct{}
	done      chan struct{}
	isRunning bool
	interval  time.Duration
	StopTime  time.Duration
	ctx       context.Context
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// NewMonitorPolicy creates a new monitor
func NewMonitorPolicy(ctx context.Context) *MonitorPolicy {
	return &MonitorPolicy{
		ctx:      ctx,
		stopChan: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (m *MonitorPolicy) Monitor(name string, fn func(), ti time.Duration) error {
	if m.IsRunning() {
		return fmt.Errorf("policy %s already running", name)
	}

	if ti == 0 {
		return fmt.Errorf("interval cannot be 0")
	}

	if fn == nil {
		return fmt.Errorf("no func to monitor")
	}

	m.interval = ti
	m.policy = name
	m.fn = fn

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.hub()
	}()

	return nil
}

func (m *MonitorPolicy) hub() error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.mu.Lock()
	m.isRunning = true
	m.mu.Unlock()

	for {
		select {
		case <-ticker.C:
			m.fn()
		case <-m.stopChan:
			return m.stop()
		case <-m.ctx.Done():
			m.mu.Lock()
			m.isRunning = false
			m.mu.Unlock()
			return nil
		}
	}
}

func (m *MonitorPolicy) Stop() {
	select {
	case m.stopChan <- struct{}{}:
	default:
	}
}

func (m *MonitorPolicy) stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	m.isRunning = false

	go func() {
		m.wg.Wait()
		close(m.done)
	}()

	t := m.StopTime
	if t == 0 {
		t = 1 * time.Second
	}

	select {
	case <-m.done:
		return nil
	case <-time.After(t):
		return fmt.Errorf("stop timeout for policy %s", m.policy)
	}
}

func (m *MonitorPolicy) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isRunning
}
