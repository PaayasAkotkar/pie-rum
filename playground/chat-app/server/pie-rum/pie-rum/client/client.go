// Package client ...
package client

import (
	"context"
	"fmt"
	rumrpc "pie-rum-sdk/misc/rum"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PieRum struct {
	conn  *grpc.ClientConn
	inner rumrpc.OnRumServiceClient
	cfg   *config
}

func New(addr string, opts ...Option) (*PieRum, error) {
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

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("xrpc: dial %s: %w", addr, err)
	}

	return &PieRum{
		conn:  conn,
		inner: rumrpc.NewOnRumServiceClient(conn),
		cfg:   &cfg,
	}, nil
}

func (c *PieRum) Close() error {
	return c.conn.Close()
}

func (c *PieRum) POST(ctx context.Context, req *rumrpc.IPostRequest) (*rumrpc.IPostResponse, error) {
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

func (r *PieRum) MonitorTag(ctx context.Context, in *rumrpc.IMonitorTagRequest) (*rumrpc.IMonitorTagResponse, error) {
	client := r.inner
	return client.MonitorTag(ctx, in)
}

func (r *PieRum) GetDoc(ctx context.Context, in *rumrpc.IDocRequest) (*rumrpc.IDocResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.GetDoc(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}

func (r *PieRum) Release(ctx context.Context, in *rumrpc.ReleaseRequest) (*rumrpc.ReleaseResponse, error) {

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

func (r *PieRum) UpdateProfileConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {

	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.UpdateProfileConfig(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}

func (r *PieRum) UpdateKitConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.UpdateKitConfig(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *PieRum) UpdateServiceConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.UpdateServiceConfig(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *PieRum) UpdateDispatcherConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.UpdateDispatcherConfig(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
func (r *PieRum) UpdateEventConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("xrpc: req must not be nil")
	}
	ctx, cancel := r.cfg.callContext(ctx)
	defer cancel()

	resp, err := r.inner.UpdateEventConfig(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("xrpc: POST: %w", err)
	}
	return resp, nil
}
