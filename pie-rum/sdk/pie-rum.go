package pierum

import (
	rum "pie-rum-sdk/pie-rum/core"

	"context"

	"google.golang.org/grpc"
)

type PieRum[In, Out any] struct {
	core *rum.PieRum[In, Out]
}
type ServerConfig struct {
	Network       string
	Address       string
	ServerOptions []grpc.ServerOption
}
type ISearch struct {
	profile, kit, service, dispatcher, event bool
}

func New[In, Out any](ctx context.Context, store *Store[In, Out]) *PieRum[In, Out] {
	return &PieRum[In, Out]{
		core: rum.New(ctx, store.core),
	}
}

// Serve starts the rum server
// tip: this can be even push in the any of the r.Go... method
func (r *PieRum[In, Out]) Serve(ctx context.Context, cnf ServerConfig) {
	r.core.Serve(ctx, rum.Server{
		Network:       cnf.Network,
		Address:       cnf.Address,
		ServerOptions: cnf.ServerOptions,
	})
}

func (r *PieRum[In, Out]) SetStore(ctx context.Context, s *Store[In, Out]) {
	r.core = rum.New(ctx, s.core)
}

func (r *PieRum[In, Out]) GoMonitor(ctx context.Context, handler rum.Handler, fns ...func()) {
	r.core.GoMonitor(ctx, handler, fns...)
}

func (r *PieRum[In, Out]) GoAuto(ctx context.Context, handler rum.Handler, fns ...func()) {
	r.core.GoAuto(ctx, handler, fns...)
}

func (r *PieRum[In, Out]) GoPoll(ctx context.Context, profile string, handler rum.Handler, fns ...func()) {
	r.core.GoPoll(ctx, profile, handler, fns...)
}

func (r *PieRum[In, Out]) GET(profile string, handler func(<-chan *rum.IResults)) {
	out := r.core.Paper(profile)
	handler(out)
}

func (r *PieRum[In, Out]) MonitorError(ctx context.Context, handler func(e error), fns ...func()) {
	r.core.GoErrors(ctx, handler, fns...)
}

// func (r*PieRum[In, Out]) Serve(ctx context.Context, cnf rum.RumServer) {
// 	r.core.Serve(ctx, cnf)
// }
