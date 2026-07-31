package server

import (
	structure "app/server/admin"
	rcrypt "app/server/crypt"
	"app/server/server/graph"
	"app/server/server/graph/model"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	pierum "pie-rum-sdk/pie-rum/core"
	rumsdk "pie-rum-sdk/pie-rum/sdk"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	coderws "github.com/coder/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func Serve() {
	const defaultPort = "8080"
	a := gin.New()
	a.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	app := New(a)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: app,
	}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GRAPHQL{})
	srv.AddTransport(transport.Websocket{
		Implementation: transport.CoderWebsocketImplementation{
			AcceptOptions: coderws.AcceptOptions{
				OriginPatterns: []string{"http://localhost:3000"},
			},
		},

		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	app.srv = srv

	//server := &http.Server{
	//	Addr: ":" + defaultPort,
	//}

	cnf := rumsdk.ServerConfig{
		Network: rcrypt.PIERUMNETWORK,
		Address: rcrypt.PIERUMADRESS,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	pie := structure.Init(ctx)

	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in pie.GoMonitor: %v", r)
			}
		}()
		pie.GoMonitor(ctx, func(result *pierum.IResults) {
			if result.IsReady {
				log.Println("res: ", result.Resuts)
				for _, r := range result.Resuts {
					var rs structure.Res
					var in structure.Req
					log.Println("input: ", string(r.Input))
					log.Println("output: ", string(r.Output))
					log.Println("dog-report", string(r.DogReport))
					log.Println("meta-report", r.MetaInfo)

					if err := json.Unmarshal(r.Input, &in); err != nil {
						log.Println("Error unmarshaling input:", err)
						continue
					}
					if err := json.Unmarshal(r.Output, &rs); err != nil {
						log.Println("Error unmarshaling output:", err)
						continue
					}

					if in.Student.ID == nil {
						log.Println("Student ID is nil")
						continue
					}

					if rs.Err != nil {
						log.Println("Error from AI model:", rs.Err)
						// Publish error to frontend
						x := rs.Err.Error()
						errorReply := &model.OnChessCoachReply{
							Information: &model.ChessCoachReply{
								Desc: &x,
							},
						}
						app.cheetah.Publish(*in.Student.ID, errorReply)
						continue
					}

					if rs.Coach != nil {
						app.cheetah.Publish(*in.Student.ID, rs.Coach)
					} else {
						log.Println("Coach result is nil")
					}
				}
			}
		})
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in pie.Serve: %v", r)
			}
		}()
		pie.Serve(ctx, cnf)
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in app.Monitor: %v", r)
			}
		}()
		app.Monitor()
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in app.Playground: %v", r)
			}
		}()
		app.Playground()
	}()

	go func() {
		defer wg.Done()
		app.Listen(":" + defaultPort)
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("🛑 Initiating graceful shutdown...")

		// Close pie-rum client connection
		log.Println("⏹️  Closing pie-rum client...")
		app.cli.Close()
		log.Println("✅ Pie-rum client closed")

		// Close pie-rum connections
		log.Println("⏹️  Closing pie-rum connections...")
		pie.Close()
		log.Println("✅ Pie-rum connections closed")

		log.Println("✅ Graceful shutdown complete")
	}()

	wg.Wait()
}
