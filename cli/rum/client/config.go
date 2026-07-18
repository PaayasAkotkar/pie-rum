package client

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

// WithAPIKey attaches an API key to every RPC as request metadata.
func WithAPIKey(key string) Option {
	return optionFunc(func(c *config) { c.apiKey = key })
}

// WithTimeout sets the per-call deadline. Defaults to 10s.
func WithTimeout(d time.Duration) Option {
	return optionFunc(func(c *config) { c.timeout = d })
}

// WithTLSConfig overrides the TLS configuration used when dialing.
func WithTLSConfig(tls *tls.Config) Option {
	return optionFunc(func(c *config) { c.tlsCfg = tls })
}

// WithInsecure disables TLS. Never use in production.
func WithInsecure() Option {
	return optionFunc(func(c *config) { c.insecure = true })
}

// WithDialOptions appends raw grpc.DialOptions for advanced use cases.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return optionFunc(func(c *config) { c.extraDial = append(c.extraDial, opts...) })
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

func (c *config) dialOptions() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	// Transport security
	switch {
	case c.insecure:
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	case c.tlsCfg != nil:
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(c.tlsCfg)))
	default:
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	}

	// API key injected per-RPC via PerRPCCredentials
	if c.apiKey != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(apiKeyCredential{
			key:      c.apiKey,
			needsTLS: !c.insecure,
		}))
	}

	// Keepalive: match what most cloud load-balancers expect
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                20 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}))

	opts = append(opts, c.extraDial...)
	return opts, nil
}

// callContext returns a child context with the configured deadline.
func (c *config) callContext(parent context.Context) (context.Context, context.CancelFunc) {
	if c.timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, c.timeout)
}

type apiKeyCredential struct {
	key      string
	needsTLS bool
}

func (a apiKeyCredential) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"x-api-key": a.key}, nil
}

func (a apiKeyCredential) RequireTransportSecurity() bool { return a.needsTLS }

type contextKey struct{}

func NewOutgoingContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 != 0 {
		panic("xrpc: NewOutgoingContext requires an even number of key/value args")
	}
	md := metadata.New(nil)
	for i := 0; i < len(kv); i += 2 {
		md.Append(kv[i], kv[i+1])
	}
	return metadata.NewOutgoingContext(ctx, md)
}
