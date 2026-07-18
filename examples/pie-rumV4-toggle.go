package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	rumrpc "pie-rum-sdk/misc/rum"
	"pie-rum-sdk/pie-rum/client"
	rumcore "pie-rum-sdk/pie-rum/core"
	rumsdk "pie-rum-sdk/pie-rum/sdk"
	"syscall"
	"time"
)

// PlayPIERUMToggleSystemExample demonstrates the toggle system (activate, deactivate, swap)
// using the rum/sdk for profiles, kits, services, dispatchers, and events
func PlayPIERUMToggleSystemExample() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	const (
		serverAddr = "localhost:9306"
	)

	// 1. Initialize Store using the SDK
	store := rumsdk.NewStore[Req, *Resp](ctx)

	store.AddProfile("profile_a", func(p *rumsdk.Profile[Req, *Resp]) {
		p.AddKit("kit_a_primary", func(kit *rumsdk.Kit[Req, *Resp]) {
			kit.AddService("service_a_v1", func(service *rumsdk.Service[Req, *Resp]) {
				service.AddDispatcher("dispatcher_a", func(dispatcher *rumsdk.Dispatcher[Req, *Resp]) {
					dispatcher.AddEvent("event_a", func(event *rumsdk.Event[Req, *Resp]) {
						event.RegisterFunc(func(ctx context.Context, req Req) (*Resp, error) {
							if req.Query == nil {
								q := "default query for profile A"
								req.Query = &q
							}
							log.Printf("[Profile A - V1] Processing request: %s", *req.Query)
							return &Resp{
								Info: fmt.Sprintf("Profile A V1: Handled request from %s with query: %s",
									*req.Name, *req.Query),
							}, nil
						})
					})
				})
			})
		})
	})

	store.AddProfile("profile_b", func(p *rumsdk.Profile[Req, *Resp]) {
		p.AddKit("kit_b_primary", func(kit *rumsdk.Kit[Req, *Resp]) {
			kit.AddService("service_b_v1", func(service *rumsdk.Service[Req, *Resp]) {
				service.AddDispatcher("dispatcher_b", func(dispatcher *rumsdk.Dispatcher[Req, *Resp]) {
					dispatcher.AddEvent("event_b", func(event *rumsdk.Event[Req, *Resp]) {
						event.RegisterFunc(func(ctx context.Context, req Req) (*Resp, error) {
							if req.Query == nil {
								q := "default query for profile B"
								req.Query = &q
							}
							log.Printf("[Profile B - V1] Processing request: %s", *req.Query)
							return &Resp{
								Info: fmt.Sprintf("Profile B V1: Handled request from %s with query: %s",
									*req.Name, *req.Query),
							}, nil
						})
					})
				})
			})
		})
	})

	store.Build()

	server := rumsdk.New(ctx, store)
	server.SetStore(ctx, store)

	log.Println("🚀 Rum V4 Toggle System Example")
	log.Printf("   Server: %s\n", serverAddr)
	log.Printf("   Profiles: profile_a, profile_b\n")

	server.GoMonitor(ctx, func(result *rumcore.IResults) {
		for _, r := range result.Resuts {
			log.Printf("✓ Monitor Output: %s", string(r.Output))
			log.Printf("✓ Monitor Input: %s", string(r.Input))
			log.Printf("✓ Monitor Report: %s", string(r.DogReport))

		}
	}, func() {
		server.Serve(ctx, rumsdk.ServerConfig{Network: "tcp", Address: serverAddr})
	},
	)

	go func() {
		time.Sleep(3 * time.Second)
		log.Println("[Client] Sending request to profile_a...")

		id := "req-001"
		name := "User_A"
		query := "Test profile A"

		req := Req{
			ID:      &id,
			Name:    &name,
			Query:   &query,
			Profile: "profile_a",
		}

		if err := sendToggleRequest(serverAddr, req); err != nil {
			log.Printf("Failed to send request: %v", err)
		}
	}()

	go func() {
		time.Sleep(6 * time.Second)
		log.Println("[Client] Sending request to profile_b...")

		id := "req-002"
		name := "User_B"
		query := "Test profile B"

		req := Req{
			ID:      &id,
			Name:    &name,
			Query:   &query,
			Profile: "profile_b",
		}

		if err := sendToggleRequest(serverAddr, req); err != nil {
			log.Printf("Failed to send request: %v", err)
		}
	}()

	go func() {
		time.Sleep(9 * time.Second)
		log.Println("[Toggle] Demonstrating profile activation/deactivation...")

		// Note: In a real implementation, you would use the gRPC client to call
		// ActivateProfile, DeactivateProfile, SwapProfile methods
		// This example shows the structure for toggle operations

		log.Println("[Toggle] Toggle system allows:")
		log.Println("  - Activate: Enable a profile/kit/service/dispatcher/event")
		log.Println("  - Deactivate: Disable a profile/kit/service/dispatcher/event")
		log.Println("  - Swap: Exchange ranks between two components")
	}()

	// Block until signal or context timeout
	<-ctx.Done()
	log.Println("shutdown")
}

// sendToggleRequest sends a request to the Rum server for toggle system testing
func sendToggleRequest(addr string, req Req) error {
	parcel, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	post := rumrpc.IPost{
		Profile: &rumrpc.ISequence{
			Profile: req.Profile,
			Input:   parcel,
		},
		Push: true,
	}

	cli, err := client.New(addr, nil)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer cli.Close()

	_, err = cli.POST(context.Background(), &rumrpc.IPostRequest{Post: []*rumrpc.IPost{&post}})
	if err != nil {
		return fmt.Errorf("post error: %w", err)
	}
	return nil
}
