package injection

import (
	"context"
	"fmt"
	cheetah "pie-rum-sdk/cheetah"
	injection "pie-rum-sdk/di/core"
	"reflect"
	"time"
)

// Client implemented client so that it can be more roboust
type Client struct {
	injection *injection.Injection
	pending   map[reflect.Type]*injection.Entry
	err       error
	cheetah   *cheetah.Cheetah[string, string]
}

// NewClient creates a new DI client
func NewClient(ctx context.Context, nodeID string) *Client {
	return &Client{
		injection: injection.New(ctx, nodeID),
		pending:   make(map[reflect.Type]*injection.Entry),
		cheetah:   cheetah.New[string, string](100),
	}
}

// AddSingleton create once and shutdown
func (c *Client) AddSingleton(req *ServiceRequest) *Client {
	if c.err != nil {
		return c
	}
	s := &injection.ServiceRegistration{
		Type:      req.Type,
		Factory:   req.Factory,
		Lifecycle: injection.Singleton,
	}

	c.injection.AddService(s)
	return c
}

// AddTransient create for every call
func (c *Client) AddTransient(req *ServiceRequest) *Client {
	if c.err != nil {
		return c
	}
	s := &injection.ServiceRegistration{
		Type:      req.Type,
		Factory:   req.Factory,
		Lifecycle: injection.Transient,
	}
	c.injection.AddService(s)

	return c
}

// BuildStatus returns if the status is ready or not
func (c *Client) BuildStatus() chan *string {
	return c.cheetah.Subscribe("build_status")
}

// CloseBuildStatus unsubscribes from the build status event via cheetah
func (c *Client) CloseBuildStatus(status chan *string) {
	c.cheetah.Unsubscribe("build_status", status)
}

// AddPooled creates a pool of instances and manages them through a pool.
func (c *Client) AddPooled(req *ServiceRequest, poolConfig *injection.PoolConfig) *Client {
	if c.err != nil {
		return c
	}
	c.injection.AddService(&injection.ServiceRegistration{
		Type:       req.Type,
		Factory:    req.Factory,
		Lifecycle:  injection.Pooled,
		PoolConfig: poolConfig,
	})

	return c
}

// Build builds the service dependency graph as per the profile
// note: make sure to subscribe to the buildstatus before calling it
// note: it can be triggered on scaleup using TriggerRebuild
func (c *Client) Build(ctx context.Context) error {
	if c.err != nil {
		return c.err
	}
	resultCh := make(chan error, 1)
	err := c.injection.BuildServices(resultCh)
	select {
	case err := <-err:
		if err == nil {
			msg := "ready"
			c.cheetah.Publish("build_status", &msg)
		}
		return err
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout building services")
	case <-ctx.Done():
		return fmt.Errorf("context cancelled")
	}
}

// GetService returns the service from the registry as per the profile
// note: make sure to use the buildstatus before calling it
func (c *Client) GetService(t reflect.Type) (any, error) {
	responseCh := make(chan *injection.ServiceResponse, 1)

	req := &injection.ServiceRequest{
		ServiceType: t,
		ResponseCh:  responseCh,
		Timeout:     5 * time.Second,
	}
	c.injection.GetService(req)
	select {
	case c.injection.SubInjection(t) <- req:
		response := <-responseCh
		return response.Instance, response.Error
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout getting service")
	}
}

// ReturnPooledService returns the service to the pool
// note: the service is returned to the pool and can be reused
func (c *Client) ReturnPooledService(t reflect.Type, instance any) error {
	return c.injection.GetContainer().ReturnPooledService(t, instance)
}

// TriggerRebuild created especailly only for the cluster scaling
func (c *Client) TriggerRebuild() error {
	if err := c.injection.RebuildSignal(); err != nil {
		return err
	}
	return nil
	// select {
	// case c.injection.rebuildSignal <- struct{}{}:
	// 	return nil
	// case <-time.After(1 * time.Second):
	// 	return fmt.Errorf("timeout triggering rebuild")
	// }
}

// Stop cleanups
func (c *Client) Stop() error {
	if c.injection == nil {
		return nil
	}
	c.Stop()

	return nil
}
