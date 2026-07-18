package injection

import (
	"context"
	"pie-rum-sdk/cheetah"
	"reflect"
)

// Injection it only does one thing is to push the channels
type Injection struct {
	container *Container

	getService    chan *ServiceRequest
	addService    chan *ServiceRegistration
	buildServices chan chan error
	rebuildSignal chan struct{}
	stopChan      chan struct{}
	done          chan struct{}

	serviceCheetah *cheetah.Cheetah[reflect.Type, ServiceRequest]
	nodeID         string
	ctx            context.Context
	cancel         context.CancelFunc
}

func New(ctx context.Context, nodeID string) *Injection {
	ctx, cancel := context.WithCancel(ctx)

	Injection := &Injection{
		serviceCheetah: cheetah.New[reflect.Type, ServiceRequest](128),
		getService:     make(chan *ServiceRequest, 100),
		addService:     make(chan *ServiceRegistration, 50),
		buildServices:  make(chan chan error, 1),
		rebuildSignal:  make(chan struct{}, 1),
		stopChan:       make(chan struct{}),
		done:           make(chan struct{}),
		nodeID:         nodeID,
		ctx:            ctx,
		cancel:         cancel,
	}

	go Injection.pipe()

	return Injection
}
