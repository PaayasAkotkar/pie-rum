package pierum

import (
	dog "pie-rum-sdk/dog/core"
	"time"
)

type IPolicy[O any] interface {
	GetName() string
	GetRank() int64
	GetConfig() *IConfig
	GetEvent(name string) *O
	GetEvents() []*O
	GetKeys() []string
	GetLen() int
	UpdateSlateChange(name string)
	UpdateSlateUsage(name string)
	SetConfig(config *IConfig)
}

type Settings struct {
	Base               time.Duration
	SleepTime          time.Duration
	Dog                dog.Settings
	EnableMetricReport bool
	MaxRequest         int64 // max request to allow
}

func (r *PieRum[In, Out]) addRequestCount() {
	r.maxRequestCount += 1
}
func (r *PieRum[In, Out]) resetRequestCount() {
	r.maxRequestCount = r.settings.MaxRequest
}

func defaultSettings() *Settings {
	return &Settings{
		EnableMetricReport: true,
		MaxRequest:         100,
	}
}

type RetryPolicy struct {
	Max      int           // how many times to retry
	Interval time.Duration // wait between retries
}

func NewRetryPolicy(max int, interval time.Duration) *RetryPolicy {
	return &RetryPolicy{
		Max:      max,
		Interval: interval,
	}
}

// set funcs

func (r *RetryPolicy) SetMaxRetry(max int) {
	r.Max = max
}
func (r *RetryPolicy) SetRetryInterval(interval time.Duration) {
	r.Interval = interval
}

// end

// get funcs

func (r *RetryPolicy) GetMaxRetry() int {
	return r.Max
}
func (r *RetryPolicy) GetRetryInterval() time.Duration {
	return r.Interval
}
