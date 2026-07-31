package injection

import (
	"fmt"
	"reflect"
	"time"
)

func (i *Injection) AddService(s *ServiceRegistration) {
	select {
	case i.addService <- s:
	case <-time.After(5 * time.Second):
		// handle timeout
	}
}

func (i *Injection) GetService(s *ServiceRequest) {
	select {
	case i.getService <- s:
	case <-time.After(5 * time.Second):
		// handle timeout
	}
}

func (i *Injection) RebuildSignal() error {
	select {
	case i.rebuildSignal <- struct{}{}:
		return nil
	case <-time.After(5 * time.Second):
		// handle timeout
		return fmt.Errorf("timeout rebuilding services")
	}
}

func (i *Injection) BuildServices(safe chan error) <-chan error {
	retCh := make(chan error, 1)

	select {
	case i.buildServices <- safe:
		return retCh
	case <-time.After(5 * time.Second):
		// handle timeout
		retCh <- fmt.Errorf("timeout building services")
		return retCh
	}
}

func (i *Injection) Stop() {
	select {
	case i.stopChan <- struct{}{}:
	case <-time.After(5 * time.Second):
		// handle timeout
	}
}

func (i *Injection) Done() {
	select {
	case i.done <- struct{}{}:
	case <-time.After(5 * time.Second):
		// handle timeout
	}
}

func (i *Injection) SubInjection(name reflect.Type) chan *ServiceRequest {
	return i.serviceCheetah.Subscribe(name)
}

func (i *Injection) UnsubInjection(name reflect.Type, sub chan *ServiceRequest) {
	i.serviceCheetah.Unsubscribe(name, sub)
}
func (i *Injection) GetContainer() *Container {
	return i.container
}
