package client

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
)

type Option interface {
	apply(*config)
}
type config struct {
	apiKey    string
	timeout   time.Duration
	tlsCfg    *tls.Config
	insecure  bool
	extraDial []grpc.DialOption
}

func defaultConfig() config {
	return config{timeout: 10 * time.Second}
}

// callContext returns a child context with the configured deadline.
func (c *config) callContext(parent context.Context) (context.Context, context.CancelFunc) {
	if c.timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, c.timeout)
}
