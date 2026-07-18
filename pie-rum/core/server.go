package pierum

import (
	"context"
	"log"
	"net"
	rumrpc "pie-rum-sdk/misc/rum"
	"runtime/debug"

	"google.golang.org/grpc"
)

type Server struct {
	Network       string
	Address       string
	ServerOptions []grpc.ServerOption
}

// Serve starts the service
func (r*PieRum[In, Out]) Serve(ctx context.Context, server Server) {
	printHeader()
	network := server.Network
	address := server.Address
	opts := server.ServerOptions
	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	rumrpc.RegisterOnRumServiceServer(grpcServer, r)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC RECOVERED in Hub goroutine: %v\n%s", r, debug.Stack())
			}
		}()
		r.Hub()
	}()

	go func() {
		<-ctx.Done()
		log.Println("shutting down gRPC server...")
		if r.DI != nil {
			r.DI.Stop()
		}
		grpcServer.GracefulStop()
	}()

	log.Println("start :)")
	if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("failed to serve: %v", err)
	}
}
