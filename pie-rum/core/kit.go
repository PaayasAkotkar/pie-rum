package pierum

import (
	"fmt"
	"pie-rum-sdk/common"
	"pie-rum-sdk/stack"
	"slices"
	"sync"
)

// const maxMetrics = 100

// IKit is the reusable container for a profile's model config, embedder,
// active/inactive services, and accumulated metrics.
type IKit[In, Out any] struct {
	mu sync.RWMutex
	// embd     *Embd
	// Model    string
	// Bucket   string
	// isHybrid bool
	// genIKit          *IGenIKit
	// sequence *ISequence[In]
	Name     string
	Rank     int64
	config   *IConfig
	registry map[string]*IService[In, Out]
	// inregistry map[string]*Service[In, Out]
	stack             *stack.Stack[string]
	slate             *ISlate
	defaultServiceKey string

	// Format          *TimeFormat
	// Metrics *IMetric `json:"metrics"`
	// Result  []*IDispatchResult
}

func NewKit[In, Out any]() *IKit[In, Out] {
	return &IKit[In, Out]{
		// Model:           model,
		// Metrics:  NewIMetric(),
		registry: make(map[string]*IService[In, Out]),
		config:   defaultConfig(),
		slate:    NewSlate(),
		stack:    stack.NewStack[string](),
		// inregistry: make(map[string]*Service[In, Out]),
	}
}

func (r *IKit[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.registry[ra]; ok {
			if rx.config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultServiceKey
}

type IResults struct {
	IsReady bool
	Resuts  []*IDispatchResult
}

func (k *IKit[In, Out]) GetName() string {
	return k.Name
}

func (k *IKit[In, Out]) GetConfig() *IConfig {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.config
}
func (k *IKit[In, Out]) GetService(key string) *IService[In, Out] {
	k.mu.RLock()
	defer k.mu.RUnlock()
	s, ok := k.registry[key]
	if !ok {
		return nil
	}
	return s
}
func (k *IKit[In, Out]) GetServices() []*IService[In, Out] {
	k.mu.RLock()
	defer k.mu.RUnlock()
	keys := k.stack.Range(k.stack.Len())
	out := make([]*IService[In, Out], 0, len(keys))
	for _, key := range keys {
		if svc, ok := k.registry[key]; ok {
			out = append(out, svc)
		}
	}
	slices.SortFunc(out, func(a, b *IService[In, Out]) int {
		return int(a.Rank - b.Rank)
	})
	return out
}
func (k *IKit[In, Out]) GetKeys() []string {
	return k.stack.Range(k.stack.Len())
}
func (k *IKit[In, Out]) GetLen() int {
	return k.stack.Len()
}
func (k *IKit[In, Out]) PushService(name string, service *IService[In, Out]) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.registry[name]; !ok {
		k.stack.Push(name)
	}
	k.registry[name] = service
}
func (k *IKit[In, Out]) RemoveService(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.registry, key)
	// delete(k.inregistry, key)
	k.stack.Erase(key)
}
func (k *IKit[In, Out]) ReplaceService(name string, service *IService[In, Out]) {
	k.registry[name] = service
}
func (k *IKit[In, Out]) SetService(services map[string]*IService[In, Out]) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.registry = services
	for key := range services {
		k.stack.Push(key)
	}
}
func (k *IKit[In, Out]) SetConfig(config *IConfig) {
	k.config = config
}

func (k *IKit[In, Out]) UpdateServiceSlateChange(name string) {
	k.slate.RecordChange(name)
}
func (k *IKit[In, Out]) UpdateServiceSlateUsage(name string) {
	k.slate.RecordUsage(name)
}

// config

func (k *IKit[In, Out]) handleServiceActivation(key string) error {
	if svc, ok := k.registry[key]; ok {

		svc.config.setActivate(true)
		k.slate.RecordChange(key)
		return nil
	}
	return nil
}

func (k *IKit[In, Out]) handleServiceDeactivation(key string) error {
	if svc, ok := k.registry[key]; ok {
		svc.config.setActivate(false)
		k.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}

func (k *IKit[In, Out]) GetRank() int64 {
	return k.Rank
}

func (k *IKit[In, Out]) SetRank(i int64) {
	k.Rank = i
}

func (k *IKit[In, Out]) handleServiceSwap(key1, key2 string) error {
	// if svc1, ok := k.registry[key1]; ok {
	// 	if svc2, ok := k.registry[key2]; ok {
	// 		swap(&svc1.Rank, &svc2.Rank)
	// 		k.slate.RecordChange(key1)
	// 		k.slate.RecordChange(key2)
	// 		return nil
	// 	}
	// 	return nil
	// }
	if err := swap(k.registry, k.registry, key1, key2); err != nil {
		return err
	}
	k.slate.RecordChange(key1)
	k.slate.RecordChange(key2)
	return nil
}

func (k *IKit[In, Out]) IsServiceActive(key string) bool {
	return k.registry[key].config.getActivate()
}

func (k *IKit[In, Out]) IsServiceSwap(key string) *IConfig {
	if seq, ok := k.registry[key]; ok {
		return seq.config
	}
	return nil
}

// end

// documentation

func (m *IKit[In, Out]) GetServicesMetadata() *IMetadata {
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
func (m *IKit[In, Out]) GetServiceMetadata(name string) *IMetadata {
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
