package pierum

import (
	"fmt"
	rumrpc "pie-rum-sdk/misc/rum"
)

type Resolved[In, Out any] struct {
	prf *IProfile[In, Out]
	kit *IKit[In, Out]
	svc *IService[In, Out]
	dt  *IDispatcher[In, Out]
}

type ActionFn[In, Out any] func(res *Resolved[In, Out], t IConfigRequest) error

type ActionEntry[In, Out any] struct {
	depth  int
	action ActionFn[In, Out]
}
type Swap interface {
	GetRank() int64
	SetRank(i int64)
}

func swap[T Swap](a, b map[string]T, akey, bkey string) error {
	if _, ok := a[akey]; !ok {
		return fmt.Errorf("key %s not found", akey)
	}
	if _, ok := b[bkey]; !ok {
		return fmt.Errorf("key %s not found", bkey)
	}
	x, y := a[akey], b[bkey]
	temp := x.GetRank()
	x.SetRank(y.GetRank())
	y.SetRank(temp)
	a[akey] = x
	b[bkey] = y
	return nil
}

func rpcToPIERUM(c *rumrpc.IConfig) *IConfig {
	var swapOverview *ISwitch
	if c.Swap != nil {
		swapOverview = &ISwitch{
			Name:    c.Swap.With,
			HSwitch: c.Swap.Swap,
		}
	}
	return &IConfig{
		Activate:     c.Activate,
		SwapOverview: swapOverview,
	}
}
