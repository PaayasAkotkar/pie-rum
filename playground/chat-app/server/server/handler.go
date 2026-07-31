package server

import (
	"app/server/server/graph/model"
	rumrpc "pie-rum-sdk/misc/rum"
)

func conGQLToRPC(co *model.IConfig) *rumrpc.IConfig {
	var swap *rumrpc.ISwap
	if co.Swap != nil {
		swap = &rumrpc.ISwap{
			Swap:    co.Swap.Swap,
			With:    co.Swap.With,
			Current: co.Swap.Current,
		}
	}
	return &rumrpc.IConfig{
		Activate: co.Activate,
		Swap:     swap,
	}
}
