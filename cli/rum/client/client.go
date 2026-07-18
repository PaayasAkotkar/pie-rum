// Package client ...
package client

import (
	"context"
	"fmt"
	rumrpc "pie-rum/misc/rum"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Rum struct {
	conn  *grpc.ClientConn
	inner rumrpc.OnRumServiceClient
	cfg   config
}

func New(addr string, opts ...Option) (*Rum, error) {
	if addr == "" {
		return nil, fmt.Errorf("xrpc: addr must not be empty")
	}

	cfg := defaultConfig()

	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(&cfg)
	}

	// dialOpts, err := cfg.dialOptions()
	// if err != nil {
	// 	return nil, fmt.Errorf("xrpc: build dial options: %w", err)
	// }

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("xrpc: dial %s: %w", addr, err)
	}

	return &Rum{
		conn:  conn,
		inner: rumrpc.NewOnRumServiceClient(conn),
		cfg:   cfg,
	}, nil
}

func (c *Rum) Close() error {
	return c.conn.Close()
}

func (c *Rum) POST(ctx context.Context, req *rumrpc.IPostRequest) (*rumrpc.IPostResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := c.cfg.callContext(ctx)
	defer cancel()

	resp, err := c.inner.POST(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) MonitorTag(ctx context.Context, in *rumrpc.IMonitorTagRequest) (*rumrpc.IMonitorTagResponse, error) {

	client := r.inner

	return client.MonitorTag(ctx, in)
}
func (r *Rum) Release(ctx context.Context, in *rumrpc.ReleaseRequest) (*rumrpc.ReleaseResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.Release(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) ActivateProfile(ctx context.Context, in *rumrpc.IActivateProfileRequest) (*rumrpc.IActivateProfileResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.ActivateProfile(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) DeactivateProfile(ctx context.Context, in *rumrpc.IDeactivateProfileRequest) (*rumrpc.IDeactivateProfileResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.DeactivateProfile(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) SwapProfile(ctx context.Context, in *rumrpc.ISwapProfileRequest) (*rumrpc.ISwapProfileResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.SwapProfile(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) DeactivateKit(ctx context.Context, in *rumrpc.IDeactivateKitRequest) (*rumrpc.IDeactivateKitResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.DeactivateKit(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) ActivateKit(ctx context.Context, in *rumrpc.IActivateKitRequest) (*rumrpc.IActivateKitResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.ActivateKit(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) SwapKit(ctx context.Context, in *rumrpc.ISwapKitRequest) (*rumrpc.ISwapKitResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.SwapKit(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) DeactivateService(ctx context.Context, in *rumrpc.IDeactivateServiceRequest) (*rumrpc.IDeactivateServiceResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.DeactivateService(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) ActivateService(ctx context.Context, in *rumrpc.IActivateServiceRequest) (*rumrpc.IActivateServiceResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.ActivateService(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) SwapService(ctx context.Context, in *rumrpc.ISwapServiceRequest) (*rumrpc.ISwapServiceResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.SwapService(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) DeactivateDispatcher(ctx context.Context, in *rumrpc.IDeactivateDispatcherRequest) (*rumrpc.IDeactivateDispatcherResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.DeactivateDispatcher(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil

}
func (r *Rum) ActivateDispatcher(ctx context.Context, in *rumrpc.IActivateDispatcherRequest) (*rumrpc.IActivateDispatcherResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.ActivateDispatcher(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) SwapDispatcher(ctx context.Context, in *rumrpc.ISwapDispatcherRequest) (*rumrpc.ISwapDispatcherResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.SwapDispatcher(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) DeactivateEvent(ctx context.Context, in *rumrpc.IDeactivateEventRequest) (*rumrpc.IDeactivateEventResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.DeactivateEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) ActivateEvent(ctx context.Context, in *rumrpc.IActivateEventRequest) (*rumrpc.IActivateEventResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.ActivateEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *Rum) SwapEvent(ctx context.Context, in *rumrpc.ISwapEventRequest) (*rumrpc.ISwapEventResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.SwapEvent(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil

}
