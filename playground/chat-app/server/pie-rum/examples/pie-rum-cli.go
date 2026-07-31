package example

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	rumcore "pie-rum-sdk/pie-rum/core"
	rumsdk "pie-rum-sdk/pie-rum/sdk"
	"syscall"
	"time"
)

// PlayPIERUMCLIBasicExample demonstrates using the SDK with GoMonitor and the client for testing
func PlayPIERUMCLIBasicExample() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	const (
		serverAddr  = "localhost:9305"
		profileName = "sdk_test_profile"
	)

	store := rumsdk.NewStore[Req, *Resp](ctx)

	store.AddProfile(profileName, func(p *rumsdk.Profile[Req, *Resp]) {
		p.AddKit("test_kit", func(kit *rumsdk.Kit[Req, *Resp]) {
			kit.AddService("test_service", func(service *rumsdk.Service[Req, *Resp]) {
				service.AddDispatcher("test_dispatcher", func(dispatcher *rumsdk.Dispatcher[Req, *Resp]) {
					dispatcher.AddEvent("test_event", func(event *rumsdk.Event[Req, *Resp]) {
						event.RegisterFunc(func(ctx context.Context, req Req) (*Resp, error) {
							if req.Query == nil {
								q := "default sdk query"
								req.Query = &q
							}

							log.Printf("[Server] Processing request for profile: %s, query: %s", req.Profile, *req.Query)

							return &Resp{
								Info: fmt.Sprintf("Success! Handled request from %s for profile %s with query: %s",
									*req.Name, req.Profile, *req.Query),
							}, nil
						})
					})
				})
			})
		})
	})

	store.Build()

	server := &rumsdk.PieRum[Req, *Resp]{}
	server.SetStore(ctx, store)

	rumServer := rumsdk.New(ctx, store)

	log.Println("🚀 Rum V4 SDK Example")
	log.Printf("   Server: %s\n", serverAddr)
	log.Printf("   Profile: %s\n", profileName)

	rumServer.GoMonitor(ctx, func(result *rumcore.IResults) {
		log.Println("result alert: ", result)
		for _, r := range result.Resuts {
			log.Printf("✓ Monitor Output: %s", string(r.Output))
			log.Printf("✓ Monitor Input: %s", string(r.Input))
			log.Printf("Monitor Report: %s", string(r.DogReport))
		}
	}, func() {
		rumServer.Serve(ctx, rumsdk.ServerConfig{Network: "tcp", Address: serverAddr})
	},
	)

	<-ctx.Done()
	log.Println("shutdown")
}
