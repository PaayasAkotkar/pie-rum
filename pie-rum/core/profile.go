// Package pierum ....
// flow:
// register -> Profile+ service -> save via name & perform modifcation
package pierum

import (
	"fmt"
	"pie-rum-sdk/common"
	rumstack "pie-rum-sdk/stack"
	"slices"
	"sync"
)

// IProfile manages toggle system
// note: rank is mandatory cause the sorting of the profile will be based on that
// if 0 is provided rank will be set in random order
type IProfile[In, Out any] struct {
	mu            sync.Mutex
	Rank          int64
	Name          string
	registry      map[string]*IKit[In, Out]
	config        *IConfig
	slate         *ISlate
	defaultKitKey string
	// inregistry map[string]map[ISequence[In]]*Kit[In, Out]
	stack rumstack.Stack[string]
}

func NewProfile[In, Out any]() *IProfile[In, Out] {
	return &IProfile[In, Out]{
		registry: make(map[string]*IKit[In, Out]),
		config:   defaultConfig(),
		slate:    NewSlate(),
		// inregistry: make(map[string]map[ISequence[In]]*Kit[In, Out]),
	}
}

func (r *IProfile[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.registry[ra]; ok {
			if rx.config.getActivate() {
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
	return r.config
}
func (r *IProfile[In, Out]) SetRank(i int64) {
	r.Rank = i
}
func (r *IProfile[In, Out]) GetKit(name string) *IKit[In, Out] {
	if kit, ok := r.registry[name]; ok {
		return kit
	}
	return nil
}
func (r *IProfile[In, Out]) GetKits() []*IKit[In, Out] {
	keys := r.stack.Range(r.stack.Len())
	out := make([]*IKit[In, Out], 0, len(keys))
	for _, k := range keys {
		if kit, ok := r.registry[k]; ok {
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
	if _, ok := r.registry[name]; !ok {
		r.stack.Push(name)
	}
	r.registry[name] = kit
}
func (r *IProfile[In, Out]) ReplaceKit(name string, kit *IKit[In, Out]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.registry[name]; ok {
		r.registry[name] = kit
	}
}
func (r *IProfile[In, Out]) RemoveKit(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.registry, key)
	// delete(r.inregistry, key)
	r.stack.Erase(key)
}
func (r *IProfile[In, Out]) SetKits(kits map[string]*IKit[In, Out]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registry = kits
	for k := range kits {
		r.stack.Push(k)
	}
}
func (r *IProfile[In, Out]) SetConfig(config *IConfig) {
	r.config = config
}

// config

func (p *IProfile[In, Out]) handleKitActivation(key string) error {
	if seq, ok := p.registry[key]; ok {
		seq.config.setActivate(true)
		p.slate.RecordChange(key)
		return nil
	}
	return activationError("")
}

func (p *IProfile[In, Out]) handleKitDeactivation(key string) error {
	if seq, ok := p.registry[key]; ok {
		seq.config.setActivate(false)
		p.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}

func (p *IProfile[In, Out]) handleKitSwap(key1, key2 string) error {
	// if seq1, ok := p.registry[key1]; ok {
	// 	if seq2, ok := p.registry[key2]; ok {
	// 		swap(&seq1.Rank, &seq2.Rank)
	// 		p.slate.RecordChange(key1)
	// 		p.slate.RecordChange(key2)
	// 		return nil
	// 	}
	// 	return nil
	// }

	if err := swap(p.registry, p.registry, key1, key2); err != nil {
		return err
	}
	p.slate.RecordChange(key1)
	p.slate.RecordChange(key2)
	return nil
}

func (p *IProfile[In, Out]) IsKitActive(key string) bool {
	return p.registry[key].config.getActivate()
}

func (p *IProfile[In, Out]) IsKitSwap(key string) *IConfig {
	if seq, ok := p.registry[key]; ok {
		return seq.config
	}
	return nil
}

// end

// documentation

func (m *IProfile[In, Out]) GetKitsMetadata() *IMetadata {
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

func (m *IProfile[In, Out]) GetKitMetadata(name string) *IMetadata {
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

func (k *IProfile[In, Out]) UpdateKitSlateChange(name string) {
	k.slate.RecordChange(name)
}
func (k *IProfile[In, Out]) UpdateKitSlateUsage(name string) {
	k.slate.RecordUsage(name)
}
