// Package post implements the grpc post client testing
package post

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	rumrpc "pie-rum/misc/rum"
	"pie-rum/rum/client"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

type IPostTest struct {
	Input   any    `yaml:"input,omitempty"`
	Profile string `yaml:"profile"`
}

func YamlPostTest() *cobra.Command {
	return &cobra.Command{
		Use:  "post-test [filename]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.ReadFile(args[0])
			if err != nil {
				panic(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*13)
			defer cancel()

			var m IPostTest
			if err := yaml.Unmarshal(data, &m); err != nil {
				panic(err)
			}
			// p, err := policy.UnPackInitPolicy(cmd)
			// if err != nil {
			// 	panic(err)
			// }
			addr := "localhost:9305"
			// cli, err := client.New(addr, nil)
			// if err != nil {
			// 	panic(err)
			// }
			var cli *client.Rum
			// var err error
			for i := 0; i < 5; i++ {
				cli, err = client.New(addr, nil)
				if err == nil {
					break
				}
				log.Printf("Waiting for server... (attempt %d/5)", i+1)
				time.Sleep(1 * time.Second)
			}
			if err != nil {
				panic("Server never came online")
			}
			b, err := json.Marshal(m.Input)
			if err != nil {
				panic(fmt.Sprintf("error in marshaling: %s", err))

			}
			post := rumrpc.IPost{
				Profile: &rumrpc.ISequence{
					Profile: m.Profile,
					Input:   b,
				},
				Push: true,
			}

			_, err = cli.POST(ctx, &rumrpc.IPostRequest{Post: []*rumrpc.IPost{&post}})
			if err != nil {
				panic(err)
			}

			log.Println("request sent successfully 🤩")

		},
	}
}

func Root() *cobra.Command {
	return YamlPostTest()
}
