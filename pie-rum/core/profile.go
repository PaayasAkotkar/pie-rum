// Package pierum ....
// flow:
// register -> Profile+ service -> save via name & perform modifcation
package pierum

import (
	"fmt"
	"log"
	rumstack "pie-rum-sdk/stack"
	"slices"
	"sync"
)

// IProfile manages toggle system
// note: rank is mandatory cause the sorting of the profile will be based on that
// if 0 is provided rank will be set in random order
type IProfile[In, Out any] struct {
	mu     sync.Mutex
	Rank   int64    `json:"rank"`
	Name   string   `json:"name"`
	Config *IConfig `json:"config"`

	Registry map[string]*IKit[In, Out] `json:"kit"`
	//slate         *ISlate                   `json:"-"`
	defaultKitKey string `json:"-"`
	// inRegistry map[string]map[ISequence[In]]*Kit[In, Out]
	stack rumstack.Stack[string] `json:"-"`
}

func NewProfile[In, Out any]() *IProfile[In, Out] {
	return &IProfile[In, Out]{
		Registry: make(map[string]*IKit[In, Out]),
		Config:   defaultConfig(),
		//slate:    NewSlate(),
		// inRegistry: make(map[string]map[ISequence[In]]*Kit[In, Out]),
	}
}

func (r *IProfile[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.Registry[ra]; ok {
			if rx.Config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultKitKey
}

func (r *IProfile[In, Out]) GetName() string {
	return r.Name
}
func (r *IProfile[In, Out]) GetRank() int64 {
	return r.Rank
}
func (r *IProfile[In, Out]) GetConfig() *IConfig {
	return r.Config
}
func (r *IProfile[In, Out]) SetRank(i int64) {
	r.Rank = i
}
func (r *IProfile[In, Out]) GetKit(name string) *IKit[In, Out] {
	if kit, ok := r.Registry[name]; ok {
		return kit
	}
	return nil
}
func (r *IProfile[In, Out]) GetKits() []*IKit[In, Out] {
	keys := r.stack.Range(r.stack.Len())
	out := make([]*IKit[In, Out], 0, len(keys))
	for _, k := range keys {
		if kit, ok := r.Registry[k]; ok {
			out = append(out, kit)
		}
	}
	slices.SortFunc(out, func(a, b *IKit[In, Out]) int {
		return int(a.Rank - b.Rank)
	})
	return out
}
func (r *IProfile[In, Out]) GetKeys() []string {
	return r.stack.Range(r.stack.Len())
}
func (r *IProfile[In, Out]) GetLen() int {
	return r.stack.Len()
}
func (r *IProfile[In, Out]) PushKit(name string, kit *IKit[In, Out]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Registry[name]; !ok {
		r.stack.Push(name)
	}
	r.Registry[name] = kit
	r.Registry[name].Name = name
}
func (r *IProfile[In, Out]) ReplaceKit(name string, kit *IKit[In, Out]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Registry[name]; ok {
		r.Registry[name] = kit
	}
}
func (r *IProfile[In, Out]) RemoveKit(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Registry, key)
	// delete(r.inRegistry, key)
	r.stack.Erase(key)
}
func (r *IProfile[In, Out]) SetKits(kits map[string]*IKit[In, Out]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Registry = kits
	for k := range kits {
		r.stack.Push(k)
	}
}
func (r *IProfile[In, Out]) SetConfig(config *IConfig) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in IProfile.SetConfig: %v", r)
		}
	}()
	r.Config = config
	return nil
}

func (r *IProfile[In, Out]) handleKitConfig(key string, config *IConfig) error {
	if r == nil {
		log.Println("ERROR: r is nil in handleKitConfig")
		return activationError("receiver is nil")
	}
	if r.Registry == nil {
		log.Println("ERROR: r.Registry is nil in handleKitConfig")
		return activationError("registry is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in handleKitConfig: %v", r)
		}
	}()
	log.Println("key: ", key)
	if config != nil {
		log.Printf("config activate: %v, swapOverview nil: %v", config.Activate, config.SwapOverview == nil)
	} else {
		log.Println("config is nil")
	}

	if seq, ok := r.Registry[key]; ok {
		log.Println("ok")
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in handleKitConfig.SetConfig: %v", r)
			}
		}()
		seq.SetConfig(config)
		r.Registry[key] = seq
		return nil
	}
	return activationError(fmt.Sprintf("cannot find %s ", key))
}

// Config

func (p *IProfile[In, Out]) handleKitActivation(key string) error {
	if seq, ok := p.Registry[key]; ok {
		seq.Config.setActivate(true)
		//p.slate.RecordChange(key)
		return nil
	}
	return activationError("")
}

func (p *IProfile[In, Out]) handleKitDeactivation(key string) error {
	if seq, ok := p.Registry[key]; ok {
		seq.Config.setActivate(false)
		//p.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}

func (p *IProfile[In, Out]) handleKitSwap(key1, key2 string) error {
	// if seq1, ok := p.Registry[key1]; ok {
	// 	if seq2, ok := p.Registry[key2]; ok {
	// 		swap(&seq1.Rank, &seq2.Rank)
	// 		p.slate.RecordChange(key1)
	// 		p.slate.RecordChange(key2)
	// 		return nil
	// 	}
	// 	return nil
	// }

	if err := swap(p.Registry, p.Registry, key1, key2); err != nil {
		return err
	}
	//p.slate.RecordChange(key1)
	//p.slate.RecordChange(key2)
	return nil
}

func (p *IProfile[In, Out]) IsKitActive(key string) bool {
	return p.Registry[key].Config.getActivate()
}

func (p *IProfile[In, Out]) IsKitSwap(key string) *IConfig {
	if seq, ok := p.Registry[key]; ok {
		return seq.Config
	}
	return nil
}

// end

// documentation
//
//func (m *IProfile[In, Out]) GetKitsMetadata() *IMetadata {
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
//
//func (m *IProfile[In, Out]) GetKitMetadata(name string) *IMetadata {
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
//
//func (k *IProfile[In, Out]) UpdateKitSlateChange(name string) {
//	k.slate.RecordChange(name)
//}
//func (k *IProfile[In, Out]) UpdateKitSlateUsage(name string) {
//	k.slate.RecordUsage(name)
//}
