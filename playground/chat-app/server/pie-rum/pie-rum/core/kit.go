package pierum

import (
	"fmt"
	"log"
	"pie-rum-sdk/stack"
	"slices"
	"sync"
)

// IKit is the reusable container for a profile's model Config, embedder,
type IKit[In, Out any] struct {
	mu       sync.RWMutex
	Rank     int64                         `json:"rank"`
	Name     string                        `json:"name"`
	Config   *IConfig                      `json:"config"`
	Registry map[string]*IService[In, Out] `json:"service"`
	stack    *stack.Stack[string]          `json:"-"`
	//slate             *ISlate                       `json:"-"`
	defaultServiceKey string `json:"-"`
}

func NewKit[In, Out any]() *IKit[In, Out] {
	return &IKit[In, Out]{
		Registry: make(map[string]*IService[In, Out]),
		Config:   defaultConfig(),
		//slate:    NewSlate(),
		stack: stack.NewStack[string](),
	}
}

func (r *IKit[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.Registry[ra]; ok {
			if rx.Config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultServiceKey
}

func (k *IKit[In, Out]) GetName() string {
	return k.Name
}

func (k *IKit[In, Out]) GetConfig() *IConfig {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.Config
}
func (k *IKit[In, Out]) GetService(key string) *IService[In, Out] {
	k.mu.RLock()
	defer k.mu.RUnlock()
	s, ok := k.Registry[key]
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
		if svc, ok := k.Registry[key]; ok {
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
	if _, ok := k.Registry[name]; !ok {
		k.stack.Push(name)
	}
	k.Registry[name] = service
}
func (k *IKit[In, Out]) RemoveService(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.Registry, key)
	// delete(k.inRegistry, key)
	k.stack.Erase(key)
}
func (k *IKit[In, Out]) ReplaceService(name string, service *IService[In, Out]) {
	k.Registry[name] = service
}
func (k *IKit[In, Out]) SetService(services map[string]*IService[In, Out]) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Registry = services
	for key := range services {
		k.stack.Push(key)
	}
}
func (k *IKit[In, Out]) SetConfig(Config *IConfig) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in IKit.SetConfig: %v", r)
		}
	}()
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Config = Config
	return nil

}
func (r *IKit[In, Out]) handleServiceConfig(key string, config *IConfig) error {
	if seq, ok := r.Registry[key]; ok {
		seq.SetConfig(config)
		r.Registry[key] = seq
		return nil
	}
	return nil
}

//	func (k *IKit[In, Out]) UpdateServiceSlateChange(name string) {
//		k.slate.RecordChange(name)
//	}
//
//	func (k *IKit[In, Out]) UpdateServiceSlateUsage(name string) {
//		k.slate.RecordUsage(name)
//	}
//
// // Config
func (k *IKit[In, Out]) handleServiceActivation(key string) error {
	if svc, ok := k.Registry[key]; ok {

		svc.Config.setActivate(true)
		//k.slate.RecordChange(key)
		return nil
	}
	return nil
}

func (k *IKit[In, Out]) handleServiceDeactivation(key string) error {
	if svc, ok := k.Registry[key]; ok {
		svc.Config.setActivate(false)
		//k.slate.RecordChange(key)
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

	if err := swap(k.Registry, k.Registry, key1, key2); err != nil {
		return err
	}
	//k.slate.RecordChange(key1)
	//k.slate.RecordChange(key2)
	return nil
}

func (k *IKit[In, Out]) IsServiceActive(key string) bool {
	return k.Registry[key].Config.getActivate()
}

func (k *IKit[In, Out]) IsServiceSwap(key string) *IConfig {
	if seq, ok := k.Registry[key]; ok {
		return seq.Config
	}
	return nil
}

// end

// documentation
//
//func (m *IKit[In, Out]) GetServicesMetadata() *IMetadata {
//	// update all the metadata
//	m.slate.metadata.Rebuild(buffers)
//
//	for n, r := range m.Registry {
//		if r.Config.getActivate() {
//			m.slate.metadata.AddActive(MetadataInfo{
//				Name:        n,
//				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//				UsageLen:    m.slate.usage[n],
//			})
//		} else {
//			m.slate.metadata.AddInActive(MetadataInfo{
//				Name:        n,
//				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//				UsageLen:    m.slate.usage[n],
//			})
//		}
//		if r.Config.SwapOverview.HSwitch {
//			m.slate.metadata.AddSwapped(MetadataInfo{
//				Name:        n,
//				LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//				UsageLen:    m.slate.usage[n],
//			})
//		}
//		m.slate.metadata.AddRanking(MetadataRankingInfo{
//			Name:     n,
//			UsageLen: m.slate.usage[n],
//		})
//
//	}
//
//	slices.SortFunc(m.slate.metadata.Rankings, func(a, b MetadataRankingInfo) int {
//		if a.UsageLen == b.UsageLen {
//			return int(a.UsageLen - b.UsageLen)
//		}
//		return int(a.UsageLen - b.UsageLen)
//	})
//
//	for i := range m.slate.metadata.Rankings {
//		m.slate.metadata.Rankings[i].Rank = int64(i + 1)
//	}
//
//	m.slate.metadata.SaveLen()
//
//	return m.slate.metadata
//}
//func (m *IKit[In, Out]) GetServiceMetadata(name string) *IMetadata {
//	// update all the metadata
//	m.slate.metadata.Rebuild(buffers)
//	r := m.Registry[name]
//	n := name
//
//	if r.Config.getActivate() {
//		m.slate.metadata.AddActive(MetadataInfo{
//			Name:        n,
//			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//			UsageLen:    m.slate.usage[n],
//		})
//	} else {
//		m.slate.metadata.AddInActive(MetadataInfo{
//			Name:        n,
//			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//			UsageLen:    m.slate.usage[n],
//		})
//	}
//	if r.Config.SwapOverview.HSwitch {
//		m.slate.metadata.AddSwapped(MetadataInfo{
//			Name:        n,
//			LastUpdated: common.FormatDateForClient(m.slate.lastUpdate[n]),
//			UsageLen:    m.slate.usage[n],
//		})
//	}
//	m.slate.metadata.AddRanking(MetadataRankingInfo{
//		Name:     n,
//		UsageLen: m.slate.usage[n],
//	})
//
//	slices.SortFunc(m.slate.metadata.Rankings, func(a, b MetadataRankingInfo) int {
//		if a.UsageLen == b.UsageLen {
//			return int(a.UsageLen - b.UsageLen)
//		}
//		return int(a.UsageLen - b.UsageLen)
//	})
//
//	for i := range m.slate.metadata.Rankings {
//		m.slate.metadata.Rankings[i].Rank = int64(i + 1)
//	}
//
//	m.slate.metadata.SaveLen()
//
//	return m.slate.metadata
//}

// end
