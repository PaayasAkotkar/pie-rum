package pierum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pie-rum-sdk/common"
	rumdog "pie-rum-sdk/dog/core"
	dog "pie-rum-sdk/dog/sdk"
	"pie-rum-sdk/stack"
	"slices"
	"sync"
	"time"

	"github.com/avast/retry-go/v5"
)

// IDispatcher controls registered agent functions and their results
type IDispatcher[in, out any] struct {
	registry map[string]*IEvent[in, out]
	events   stack.Stack[string]
	Settings Settings
	Rank     int64
	Name     string
	config   *IConfig
	slate    *ISlate
	stack    *stack.Stack[string]
	// mutable
	// rinput     map[string]in
	result     map[string]*IDispatchResult
	metric     map[string]map[int]IAgentResp // name -> count -> resp
	isComplete map[string]bool
	// end
	defaultEventKey string

	wg sync.WaitGroup
}

// IAgentResp holds per-call metric data
type IAgentResp struct {
	Succeed *IMetricAgentSucceed `json:"succeed"`
	Fail    *IMetricAgentFail    `json:"fail"`
}

func NewDispatcher[in, out any](settings Settings) *IDispatcher[in, out] {
	return &IDispatcher[in, out]{
		registry: make(map[string]*IEvent[in, out]),
		// rinput:     make(map[string]in),
		stack:      stack.NewStack[string](),
		Settings:   settings,
		slate:      NewSlate(),
		config:     defaultConfig(),
		result:     make(map[string]*IDispatchResult),
		isComplete: make(map[string]bool),
		metric:     make(map[string]map[int]IAgentResp),
	}
}

func (r *IDispatcher[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.registry[ra]; ok {
			if rx.config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultEventKey
}

func (d *IDispatcher[in, out]) release() {

	// for r := range d.rinput {
	// 	delete(d.rinput, r)
	// }
	for r := range d.result {
		delete(d.result, r)
	}
	for r := range d.isComplete {
		delete(d.isComplete, r)
	}
	for r := range d.metric {
		delete(d.metric, r)
	}
}

func (d *IDispatcher[in, out]) GetName() string {
	return d.Name
}
func (d *IDispatcher[in, out]) GetRank() int64 {
	return d.Rank
}
func (d *IDispatcher[in, out]) GetConfig() *IConfig {
	return d.config
}
func (d *IDispatcher[in, out]) GetEvent(name string) *IEvent[in, out] {
	return d.registry[name]
}
func (d *IDispatcher[In, Out]) GetEvents() []*IEvent[In, Out] {
	keys := d.events.Range(d.events.Len())
	out := make([]*IEvent[In, Out], 0, len(keys))
	for _, key := range keys {
		if svc, ok := d.registry[key]; ok {
			out = append(out, svc)
		}
	}
	slices.SortFunc(out, func(a, b *IEvent[In, Out]) int {
		return int(a.Rank - b.Rank)
	})
	return out
}
func (d *IDispatcher[in, out]) GetKeys() []string {
	return d.events.Max()
}
func (d *IDispatcher[in, out]) GetLen() int {
	return d.stack.Len()
}

func (k *IDispatcher[In, Out]) UpdateEventSlateChange(name string) {
	k.slate.RecordChange(name)
}
func (k *IDispatcher[In, Out]) UpdateEventSlateUsage(name string) {
	k.slate.RecordUsage(name)
}

func (d *IDispatcher[in, out]) PushEvent(event string, fn *IEvent[in, out]) {
	if _, ok := d.registry[event]; !ok {
		d.events.Push(event)
	}
	d.registry[event] = fn
}
func (d *IDispatcher[In, Out]) ReplaceEvent(name string, fn *IEvent[In, Out]) {
	if _, ok := d.registry[name]; ok {
		d.registry[name] = fn
	}
}
func (d *IDispatcher[in, out]) RemoveEvent(name string) {
	delete(d.registry, name)
	// delete(d.rinput, name)
	delete(d.isComplete, name)
	delete(d.metric, name)
	delete(d.result, name)
	d.events.Erase(name)
}
func (d *IDispatcher[in, out]) SetEvents(events map[string]*IEvent[in, out]) {
	d.registry = events
	for v := range events {
		if _, ok := d.registry[v]; !ok {
			d.events.Push(v)
		}
	}
}
func (d *IDispatcher[in, out]) SetConfig(config *IConfig) {
	d.config = config
}

// get funcs

// func (d *IDispatcher[in, out]) GetRegistry() map[string]*IEvent[in, out] {
// 	return d.registry
// }

func (d *IDispatcher[in, out]) GetResults(name string) *IDispatchResult {
	if _, ok := d.result[name]; !ok {
		for n := range d.registry {
			log.Println("names: ", n)
		}
		log.Println("IDispatcher: not found name ", name)
		return nil
	}
	return d.result[name]
}

// GetMetric returns the latest metric entry for a named dispatch
func (d *IDispatcher[in, out]) GetMetric(name string) IAgentResp {
	return d.metric[name][d.metricCount(name)]
}

func (d *IDispatcher[in, out]) GetMetrics(name string) map[int]IAgentResp {
	return d.metric[name]
}

func (d *IDispatcher[in, out]) metricCount(name string) int {
	return len(d.metric[name])
}

// end

// set funcs

// normalCall runs every registered event in rank order without caling the metric.
func (d *IDispatcher[in, out]) normalCall(ctx context.Context, input in) {
	log.Println("[normal call]")

	for name, fn := range d.registry {

		outp, err := fn.Fn(ctx, input)
		if err != nil {
			log.Println("failed", err)
			continue
		}
		m, err := json.Marshal(outp)
		if err != nil {
			log.Println("failed", err)
			continue
		}

		n, err := json.Marshal(input)
		if err != nil {
			log.Println("failed", err)
			continue
		}
		res := NewDispatchResult()
		res.Input = n
		res.Output = m
		res.IsReady = true
		d.handleOutput(name, res)
		d.UpdateEventSlateUsage(name)
		d.handleComplete(name, true)
	}
}

// metric runs every registered event in rank order.
// calls the methods and runs the metric
func (d *IDispatcher[in, out]) metricCall(ctx context.Context, input in) map[string]error {
	errs := make(map[string]error)

	for name, fn := range d.registry { // rank-sorted, all events, never skipped

		cfg := fn.getConfig()
		if cfg != nil {
			if !cfg.getActivate() {
				err := fmt.Errorf("event %s is deactivated", name)
				d.writeMetric(name, IAgentResp{Fail: &IMetricAgentFail{
					At:     time.Now(),
					Reason: err.Error(),
					Type:   "deactive",
				}})
				errs[name] = err
				continue // record & move on, never break
			}

			if sw := cfg.GetSwapOverview(); sw != nil && sw.HSwitch && sw.Name != "" {
				swFn, swOk := d.registry[sw.Name]
				if !swOk {
					err := fmt.Errorf("swapped event %s (from %s) not found", sw.Name, name)
					d.writeMetric(name, IAgentResp{Fail: &IMetricAgentFail{
						At:     time.Now(),
						Reason: err.Error(),
						Type:   "swap",
					}})
					errs[name] = err
					continue
				}
				swCfg := swFn.getConfig()
				if swCfg != nil && !swCfg.getActivate() {
					err := fmt.Errorf("swapped event %s (from %s) is deactivated", sw.Name, name)
					d.writeMetric(name, IAgentResp{Fail: &IMetricAgentFail{
						At:     time.Now(),
						Reason: err.Error(),
						Type:   "deactive",
					}})
					errs[name] = err
					continue
				}
				fn = swFn
				name = sw.Name // ← update name to the swapped target
			}
		}

		max := 1
		interval := time.Duration(0)
		if policy := fn.Retry; policy != nil {
			max = policy.Max + 1
			interval = policy.Interval
		}

		t := d.Settings.Base
		if t == 0 {
			t = 10 * time.Second
		}
		ts := d.Settings.SleepTime
		if ts == 0 {
			ts = 100 * time.Millisecond
		}

		// Capture loop variables for the closure — critical in Go loops
		capturedFn := fn
		capturedName := name

		p := rumdog.NewPolicy[out](t)
		p.Name = capturedName
		a := func() (*out, error) {
			time.Sleep(ts)
			resp, err := capturedFn.Fn(ctx, input)
			return &resp, err
		}
		p.AddFunc(rumdog.Funcs[out]{Name: capturedName, Fn: &a})

		attempt := 0
		ret := retry.New(retry.Attempts(uint(max)), retry.Delay(interval))

		err := ret.Do(func() error {
			attempt++
			if ctx.Err() != nil {
				return retry.Unrecoverable(ctx.Err())
			}

			cli := dog.NewClient[out](t)
			defer cli.Close()
			cli.DefinePolicy(capturedName, ts).AddFuncWithReturn(capturedName, a).Build()

			rep, err := cli.ExecuteAndReport(p.Name)
			if err == nil {
				inData, _ := json.Marshal(input)
				d.writeMetric(capturedName, IAgentResp{
					Succeed: &IMetricAgentSucceed{
						TimeTaken:     rep.TotalDuration,
						AgentReply:    string(rep.Output),
						ClientRequest: string(inData),
					},
				})
				r := NewDispatchResult()
				r.IsReady = true
				// r.Metric = NewProfileMetric()
				r.DogReport = rep.JSON()
				r.Output = rep.Output
				r.Input, _ = json.Marshal(input)
				d.handleOutput(capturedName, r)
				d.handleComplete(capturedName, true)
				d.UpdateEventSlateUsage(capturedName)
				cli.Unregister(p.Name)
				return nil
			}

			cli.Unregister(p.Name)
			d.writeMetric(capturedName, IAgentResp{Fail: &IMetricAgentFail{
				Reason: fmt.Sprintf("attempt %d: %s", attempt, err.Error()),
				At:     time.Now(),
			}})
			return err
		})

		if err != nil {
			errs[capturedName] = err // collect, never break
		}
	}

	return errs // empty == all events succeeded
}

func (d *IDispatcher[in, out]) writeMetric(name string, resp IAgentResp) {
	if _, ok := d.metric[name]; !ok {
		d.metric[name] = make(map[int]IAgentResp)
	}
	d.metric[name][d.metricCount(name)+1] = resp
}

func (d *IDispatcher[in, out]) handleOutput(name string, res *IDispatchResult) {
	d.result[name] = res
}

func (d *IDispatcher[in, out]) handleComplete(name string, complete bool) {
	d.isComplete[name] = complete
}

// config

func (d *IDispatcher[In, Out]) handleEventActivation(key string) error {
	if fn, ok := d.registry[key]; ok {
		fn.config.setActivate(true)
		d.slate.RecordChange(key)
		return nil
	}
	return activationError("")
}

func (d *IDispatcher[In, Out]) handleEventDeactivation(key string) error {
	if fn, ok := d.registry[key]; ok {
		fn.config.setActivate(false)
		d.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}

func (d *IDispatcher[In, Out]) SetRank(i int64) {
	d.Rank = i
}
func (d *IDispatcher[In, Out]) handleEventSwap(key1, key2 string) error {
	if err := swap(d.registry, d.registry, key1, key2); err != nil {
		return err
	}
	d.slate.RecordChange(key1)
	d.slate.RecordChange(key2)
	return nil
}

func (d *IDispatcher[in, out]) IsEventActive(key string) bool {
	return d.registry[key].config.getActivate()
}

func (d *IDispatcher[in, out]) IsEventSwap(key string) *IConfig {
	if seq, ok := d.registry[key]; ok {
		return seq.config
	}
	return nil
}

// end

// documentation

func (m *IDispatcher[In, Out]) GetEventsMetadata() *IMetadata {
	// update all the metadata
	m.slate.metadata.Rebuild(buffers)

	for n, r := range m.registry {
		if r.config.getActivate() {
			m.slate.metadata.AddActive(MetadataInfo{
				Name:        n,
				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
				UsageLen:    m.slate.usage[n],
			})
		} else {
			m.slate.metadata.AddInActive(MetadataInfo{
				Name:        n,
				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
				UsageLen:    m.slate.usage[n],
			})
		}
		if r.config.swapOverview.HSwitch {
			m.slate.metadata.AddSwapped(MetadataInfo{
				Name:        n,
				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
				UsageLen:    m.slate.usage[n],
			})
		}
		m.slate.metadata.AddRanking(MetadataRankingInfo{
			Name:     n,
			UsageLen: m.slate.usage[n],
		})

	}

	slices.SortFunc(m.slate.metadata.Rankings, func(a, b MetadataRankingInfo) int {
		if a.UsageLen == b.UsageLen {
			return int(a.UsageLen - b.UsageLen)
		}
		return int(a.UsageLen - b.UsageLen)
	})

	for i := range m.slate.metadata.Rankings {
		m.slate.metadata.Rankings[i].Rank = int64(i + 1)
	}

	m.slate.metadata.SaveLen()

	return m.slate.metadata
}
func (m *IDispatcher[In, Out]) GetEventMetadata(name string) *IMetadata {
	// update all the metadata
	m.slate.metadata.Rebuild(buffers)
	r := m.registry[name]
	n := name

	if r.config.getActivate() {
		m.slate.metadata.AddActive(MetadataInfo{
			Name:        n,
			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
			UsageLen:    m.slate.usage[n],
		})
	} else {
		m.slate.metadata.AddInActive(MetadataInfo{
			Name:        n,
			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
			UsageLen:    m.slate.usage[n],
		})
	}
	if r.config.swapOverview.HSwitch {
		m.slate.metadata.AddSwapped(MetadataInfo{
			Name:        n,
			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
			UsageLen:    m.slate.usage[n],
		})
	}
	m.slate.metadata.AddRanking(MetadataRankingInfo{
		Name:     n,
		UsageLen: m.slate.usage[n],
	})

	slices.SortFunc(m.slate.metadata.Rankings, func(a, b MetadataRankingInfo) int {
		if a.UsageLen == b.UsageLen {
			return int(a.UsageLen - b.UsageLen)
		}
		return int(a.UsageLen - b.UsageLen)
	})

	for i := range m.slate.metadata.Rankings {
		m.slate.metadata.Rankings[i].Rank = int64(i + 1)
	}

	m.slate.metadata.SaveLen()

	return m.slate.metadata
}

// end
