package pierum

import "context"

// IEvent is a callable agent function
type IEvent[in, out any] struct {
	Fn func(ctx context.Context, req in) (out, error)
	// on hold Void  * func(ctx context.Context, req in) // new
	Rank   int64 // new
	Retry  *RetryPolicy
	config *IConfig
}

func NewRegisterFunc[In, Out any]() *IEvent[In, Out] {
	return &IEvent[In, Out]{
		config: defaultConfig(),
	}
}
func (t *IEvent[in, out]) GetRank() int64 {
	return t.Rank
}

func (t *IEvent[in, out]) SetRank(i int64) {
	t.Rank = i
}
func (t *IEvent[in, out]) setConfig(c *IConfig) {
	t.config = c
}
func (t *IEvent[in, out]) getConfig() *IConfig {
	return t.config
}
func (t *IEvent[in, out]) SetRetry(r *RetryPolicy) {
	t.Retry = r
}
func (t *IEvent[in, out]) GetRetry() *RetryPolicy {
	return t.Retry
}

// Handler is the func the caller provides to handle each result
type Handler func(result *IResults)

// IDispatchResult holds the result of a completed dispatch call
type IDispatchResult struct {
	IsReady bool
	// Metric    *ProfileMetric
	DogReport []byte
	Output    []byte
	Input     []byte
}

func NewDispatchResult() *IDispatchResult {
	return &IDispatchResult{
		IsReady: false,
		// Metric:  NewProfileMetric(),
	}
}
