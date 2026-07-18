package pierum

import (
	"context"
	cheetah "pie-rum-sdk/cheetah"
	injection "pie-rum-sdk/di/sdk"
	rumrpc "pie-rum-sdk/misc/rum"
	"strings"
	"sync"
)

const (
	buffers = 1000
)

// PieRum implements the core design of server
type PieRum[In, Out any] struct {
	rumrpc.UnimplementedOnRumServiceServer

	DI *injection.Client

	cheetah *cheetah.Cheetah[string, IResults]
	store   *IStore[In, Out]

	post        chan ILinks[In, Out]
	monitorTags chan []string

	release         chan bool
	maxRequestCount int64

	activateProfile    chan []IConfigRequest
	activateKit        chan []IConfigRequest
	activateService    chan []IConfigRequest
	activateDispatcher chan []IConfigRequest
	activateEvent      chan []IConfigRequest

	deactivateProfile    chan []IConfigRequest
	deactivateKit        chan []IConfigRequest
	deactivateService    chan []IConfigRequest
	deactivateDispatcher chan []IConfigRequest
	deactivateEvent      chan []IConfigRequest
	settings             *Settings
	swapProfile          chan []IConfigRequest
	swapService          chan []IConfigRequest
	swapKit              chan []IConfigRequest
	swapDispatcher       chan []IConfigRequest
	swapEvent            chan []IConfigRequest
	cheetahDetector      *cheetah.Cheetah[string, error]
	actions              map[string]ActionEntry[In, Out]
	ctx                  context.Context
	mu                   sync.Mutex
	wg                   sync.WaitGroup
	sb                   strings.Builder
}

func New[In, Out any](ctx context.Context, store *IStore[In, Out]) *PieRum[In, Out] {

	r := &PieRum[In, Out]{
		store:              store,
		cheetah:            cheetah.New[string, IResults](1),
		cheetahDetector:    cheetah.New[string, error](1),
		monitorTags:        make(chan []string, buffers),
		release:            make(chan bool, buffers),
		post:               make(chan ILinks[In, Out]),
		maxRequestCount:    0,
		activateProfile:    make(chan []IConfigRequest, buffers),
		activateKit:        make(chan []IConfigRequest, buffers),
		activateService:    make(chan []IConfigRequest, buffers),
		activateDispatcher: make(chan []IConfigRequest, buffers),
		activateEvent:      make(chan []IConfigRequest, buffers),

		deactivateProfile:    make(chan []IConfigRequest, buffers),
		deactivateKit:        make(chan []IConfigRequest, buffers),
		deactivateService:    make(chan []IConfigRequest, buffers),
		deactivateDispatcher: make(chan []IConfigRequest, buffers),
		deactivateEvent:      make(chan []IConfigRequest, buffers),
		settings:             defaultSettings(),
		swapProfile:          make(chan []IConfigRequest, buffers),
		swapService:          make(chan []IConfigRequest, buffers),
		swapKit:              make(chan []IConfigRequest, buffers),
		swapDispatcher:       make(chan []IConfigRequest, buffers),
		swapEvent:            make(chan []IConfigRequest, buffers),
		actions:              make(map[string]ActionEntry[In, Out], buffers),

		ctx: ctx,
		DI:  injection.NewClient(ctx, "rum-server"),
	}

	r.actions = map[string]ActionEntry[In, Out]{
		activateAction + depthOne: {1, func(res *Resolved[In, Out], t IConfigRequest) error {
			return r.store.handleProfileActivation(t.Profile)
		}},
		activateAction + depthTwo: {2, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.prf.handleKitActivation(t.Kit)
		}},
		activateAction + depthThree: {3, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.kit.handleServiceActivation(t.Service)
		}},
		activateAction + depthFour: {4, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.svc.handleDispatcherActivation(t.Dispatcher)
		}},
		activateAction + depthFive: {5, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.dt.handleEventActivation(t.Event)
		}},
		deactiveAction + depthOne: {1, func(res *Resolved[In, Out], t IConfigRequest) error {
			return r.store.handleProfileDeactivation(t.Profile)
		}},
		deactiveAction + depthTwo: {2, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.prf.handleKitDeactivation(t.Kit)
		}},
		deactiveAction + depthThree: {3, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.kit.handleServiceDeactivation(t.Service)
		}},
		deactiveAction + depthFour: {4, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.svc.handleDispatcherDeactivation(t.Dispatcher)
		}},
		deactiveAction + depthFive: {5, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.dt.handleEventDeactivation(t.Event)
		}},

		swapAction + depthOne: {1, func(res *Resolved[In, Out], t IConfigRequest) error {
			return r.store.handleProfileSwap(t.Profile, t.Swap)
		}},
		swapAction + depthTwo: {2, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.prf.handleKitSwap(t.Kit, t.Swap)
		}},
		swapAction + depthThree: {3, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.kit.handleServiceSwap(t.Service, t.Swap)
		}},
		swapAction + depthFour: {4, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.svc.handleDispatcherSwap(t.Dispatcher, t.Swap)
		}},
		swapAction + depthFive: {5, func(res *Resolved[In, Out], t IConfigRequest) error {
			return res.dt.handleEventSwap(t.Event, t.Swap)
		}},
	}
	return r
}
