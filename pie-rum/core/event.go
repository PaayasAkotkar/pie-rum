package pierum

import "context"

// IEvent is a callable agent function
type IEvent[in, out any] struct {
	Fn func(ctx context.Context, req in) (out, error) `json:"-"`
	// on hold Void  * func(ctx context.Context, req in) // new
	Rank   int64        `json:"rank"` // new
	Retry  *RetryPolicy `json:"-"`
	Config *IConfig     `json:"config"`
}

func NewRegisterFunc[In, Out any]() *IEvent[In, Out] {
	return &IEvent[In, Out]{
		Config: defaultConfig(),
	}
}
func (t *IEvent[in, out]) GetRank() int64 {
	return t.Rank
}

func (t *IEvent[in, out]) SetRank(i int64) {
	t.Rank = i
}
func (t *IEvent[in, out]) setConfig(c *IConfig) {
	defer func() {
		if r := recover(); r != nil {
			// Log panic but don't crash - event config is less critical
		}
	}()
	t.Config = c
}
func (t *IEvent[in, out]) getConfig() *IConfig {
	return t.Config
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
	IsReady bool `json:"is_ready"`
	// Metric    *ProfileMetric
	DogReport []byte `json:"dog_report"`
	Output    []byte `json:"output"`
	Input     []byte `json:"input"`
	CreatedAt string `json:"created_at"`
	// this is created cause i do belive the stat of plugin packge may overload
	// the backend server
	MetaInfo IDispatchResultMetaData `json:"meta_info"` // info ingredients required to complete that process
}

type IDispatchResultMetaData struct {
	Profile     []string `json:"profiles"` // name of profiles if any
	Kits        []string `json:"kits"`
	Services    []string `json:"services"`
	Dispatchers []string `json:"dispatchers"`
	Events      []string `json:"events"`
}

func NewDispatchResult() *IDispatchResult {
	return &IDispatchResult{
		IsReady: false,
		MetaInfo: IDispatchResultMetaData{
			Profile:     make([]string, 0),
			Kits:        make([]string, 0),
			Services:    make([]string, 0),
			Dispatchers: make([]string, 0),
			Events:      make([]string, 0),
		},
		// Metric:  NewProfileMetric(),
	}
}
