package dog

import (
	"context"
	"fmt"
	"log"
	cheetah "pie-rum-sdk/cheetah"
	rumpaint "pie-rum-sdk/paint"
	"strings"
	"sync"
	"time"
)

// Dog provides a robust timeout management system
type Dog[T any] struct {
	mu sync.RWMutex

	// Core configuration
	base time.Duration

	cheetah *cheetah.Cheetah[string, WatchdogReport]

	// Policy management
	policy    map[string]*Policy[T]       // name -> policy
	lifecycle map[string]*PolicyLifecycle // name -> state tracker

	// Duration tracking per policy
	durations map[string]time.Duration

	// Progress per policy
	progress map[string]*ExeProgress

	// Health per policy
	health map[string]*Health

	// System metrics per policy
	metrics map[string]*SystemMetrics

	// Report storage
	reports map[string]*WatchdogReport

	// Channels
	register     chan *Policy[T]
	unregister   chan string
	parkDog      chan string
	done         chan IDone
	bark         chan IBark
	reset        chan string
	resetAll     chan bool
	stopCh       chan struct{}
	doneCh       chan struct{}
	summonCh     chan string
	registeredCh chan string // Signals when registration is complete

	// Monitors
	monitors map[string]*MonitorPolicy

	// Settings
	Settings *Settings

	// Control & Context
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once
}

// New creates a new Dog instance with base timeout
func New[T any](base time.Duration) *Dog[T] {
	return NewWithContext[T](context.Background(), base)
}

// NewWithContext creates a new Dog with context support
func NewWithContext[T any](ctx context.Context, base time.Duration) *Dog[T] {
	if base == 0 {
		panic("base timeout cannot be 0")
	}

	ctx, cancel := context.WithCancel(ctx)

	rd := &Dog[T]{
		base:         base,
		policy:       make(map[string]*Policy[T]),
		lifecycle:    make(map[string]*PolicyLifecycle),
		durations:    make(map[string]time.Duration),
		progress:     make(map[string]*ExeProgress),
		health:       make(map[string]*Health),
		metrics:      make(map[string]*SystemMetrics),
		reports:      make(map[string]*WatchdogReport),
		monitors:     make(map[string]*MonitorPolicy),
		register:     make(chan *Policy[T], 100),
		unregister:   make(chan string, 100),
		parkDog:      make(chan string, 100),
		done:         make(chan IDone, 1000),
		bark:         make(chan IBark, 1000),
		reset:        make(chan string, 100),
		resetAll:     make(chan bool, 10),
		stopCh:       make(chan struct{}),
		doneCh:       make(chan struct{}),
		summonCh:     make(chan string, 100),
		registeredCh: make(chan string, 100),
		cheetah:      cheetah.New[string, WatchdogReport](100),
		Settings:     DefaultSettings(),
		ctx:          ctx,
		cancel:       cancel,
	}

	printHeader()
	return rd
}

// registerPolicy adds a new policy
func (rd *Dog[T]) registerPolicy(policy *Policy[T]) {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	rd.policy[policy.Name] = policy
	rd.progress[policy.Name] = NewProgress()
	rd.health[policy.Name] = &Health{}

	title := "Registration Succeed 😄"
	desc := []string{
		fmt.Sprintf("Name: %s", policy.Name),
		fmt.Sprintf("Funcs to track: %d", len(policy.GetFunc())),
	}

	for _, fn := range policy.GetFunc() {
		desc = append(desc, fmt.Sprintf("- %s (rank: %d)", fn.Name, fn.Rank))
	}
	t := rumpaint.Card(title, strings.Join(desc, ", "))
	log.Println(t)
}

// unregisterPolicy removes a policy
func (rd *Dog[T]) unregisterPolicy(name string) {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	// if stopChan, exists := rd.tickers[name]; exists {
	// 	close(stopChan)
	// 	delete(rd.tickers, name)
	// }

	if _, exists := rd.policy[name]; exists {
		delete(rd.policy, name)
		delete(rd.progress, name)
		delete(rd.health, name)
		delete(rd.reports, name)
		fmt.Printf("[Unregister] Policy '%s' cleaned up\n", name)
	}
}

// resetPolicy releases the call counts for specific policy
func (rd *Dog[T]) resetPolicy(name string) {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	if policy, exists := rd.policy[name]; exists {
		policy.Release()
		policy.Succeed.Release()
		policy.Fail.Release()

		rd.progress[name] = NewProgress()
		rd.reports[name] = &WatchdogReport{
			PolicyName:     name,
			StartTime:      time.Now(),
			TimeLimit:      policy.GetBase(),
			FailureReasons: make([]string, 0),
		}
		fmt.Printf("[Reset] Policy '%s' reset\n", name)
	}
}

// resetPolicy releases all the call counts from all the policies
func (rd *Dog[T]) resetAllPolicies() {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	for name, policy := range rd.policy {
		policy.Release()
		policy.Succeed.Release()
		policy.Fail.Release()

		rd.progress[name] = NewProgress()
		fmt.Printf("[Reset] Policy '%s' reset\n", name)
	}
}

// getTimeFloat returns the conversion based on settings
func (rd *Dog[T]) getTimeFloat(t time.Duration) IConv {
	switch true {
	case rd.Settings.ConvDurationInHours:
		return durationToHour(t)
	case rd.Settings.ConvDurationInMins:
		return durationToMin(t)
	case rd.Settings.ConvDurationInSecs:
		return durationToSec(t)
	case rd.Settings.ConvDurationInMiliSecs:
		return durationToMS(t)
	}
	return IConv{Conv: 0, Unit: "nil"}
}
