package pierum

import (
	"context"
	"fmt"
	"pie-rum-sdk/common"
	"pie-rum-sdk/stack"
	"slices"
)

// IService the implementation guideline for the kits
// careful: only use that much func as much require
// becuase the formula will run all the dsipatch methods first than move on to the next call
// Rank current Rank of the serivce
type IService[in, out any] struct {
	context              context.Context
	Rank                 int64
	Name                 string
	slate                *ISlate
	config               *IConfig
	defaultDispatcherKey string
	registry             map[string]*IDispatcher[in, out]
	stack                *stack.Stack[string]
	// Budget   *Budget
}

func NewService[in, out any](ctx context.Context, settings Settings, name string) *IService[in, out] {
	// b := NewBudget(0, 0)
	return &IService[in, out]{
		context: ctx,
		// Format:  NewTimeFormat(),
		// Budget:   b,
		registry: make(map[string]*IDispatcher[in, out], buffers),
		slate:    NewSlate(),
		stack:    stack.NewStack[string](),
		config:   defaultConfig(),
		Rank:     1,
		Name:     name,
	}
}

func (r *IService[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.registry[ra]; ok {
			if rx.config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultDispatcherKey
}

func (d *IService[in, out]) GetName() string {
	return d.Name
}
func (d *IService[in, out]) GetRank() int64 {
	return d.Rank
}
func (d *IService[in, out]) GetConfig() *IConfig {
	return d.config
}
func (d *IService[in, out]) GetDispatcher(key string) *IDispatcher[in, out] {
	if d.registry == nil {
		return nil
	}
	return d.registry[key]
}
func (r *IService[In, Out]) GetDispatchers() []*IDispatcher[In, Out] {
	keys := r.stack.Range(r.stack.Len())
	out := make([]*IDispatcher[In, Out], 0, len(keys))
	for _, k := range keys {
		if dp, ok := r.registry[k]; ok {
			out = append(out, dp)
		}
	}
	slices.SortFunc(out, func(a, b *IDispatcher[In, Out]) int {
		return int(a.Rank - b.Rank)
	})
	return out
}
func (s *IService[in, out]) GetKeys() []string {
	return s.stack.Max()
}
func (s *IService[in, out]) GetLen() int {
	return s.stack.Len()
}
func (d *IService[in, out]) PushDispatcher(name string, dp *IDispatcher[in, out]) {
	if _, ok := d.registry[name]; !ok {
		d.stack.Push(name)
	}
	d.registry[name] = dp
}
func (s *IService[In, Out]) ReplaceDispatcher(key string, dp *IDispatcher[In, Out]) {
	if _, ok := s.registry[key]; ok {
		s.registry[key] = dp
	}
}
func (d *IService[in, out]) RemoveDispatcher(name string) {
	delete(d.registry, name)
	d.stack.Erase(name)
}
func (d *IService[in, out]) SetDispatcher(dp map[string]*IDispatcher[in, out]) {
	d.registry = dp
	for v := range dp {
		if _, ok := d.registry[v]; !ok {
			d.stack.Push(v)
		}
	}
}
func (d *IService[in, out]) SetConfig(c *IConfig) {
	d.config = c
}

func (k *IService[In, Out]) UpdateDispatcherSlateChange(name string) {
	k.slate.RecordChange(name)
}
func (k *IService[In, Out]) UpdateDispatcherSlateUsage(name string) {
	k.slate.RecordUsage(name)
}

// config

func (s *IService[in, out]) handleDispatcherActivation(key string) error {
	if d, ok := s.registry[key]; ok {
		d.config.setActivate(true)
		s.slate.RecordChange(key)
		return nil
	}
	return nil
}

func (s *IService[in, out]) handleDispatcherDeactivation(key string) error {
	if d, ok := s.registry[key]; ok {
		d.config.setActivate(false)
		s.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}
func (s *IService[in, out]) SetRank(i int64) {
	s.Rank = i
}
func (s *IService[in, out]) handleDispatcherSwap(key1, key2 string) error {
	// if d1, ok := s.registry[key1]; ok {
	// 	if d2, ok := s.registry[key2]; ok {
	// 		swap(&d1.Rank, &d2.Rank)
	// 		s.slate.RecordChange(key1)
	// 		s.slate.RecordChange(key2)
	// 		return nil
	// 	}
	// 	return nil
	// }
	if err := swap(s.registry, s.registry, key1, key2); err != nil {
		return err
	}
	s.slate.RecordChange(key1)
	s.slate.RecordChange(key2)
	return nil
}

func (p *IService[in, out]) IsDispatcherActive(key string) bool {
	return p.registry[key].config.getActivate()
}

func (p *IService[in, out]) IsDispatcherSwap(key string) *IConfig {
	if seq, ok := p.registry[key]; ok {
		return seq.config
	}
	return nil
}

// end

// documentation

func (m *IService[In, Out]) GetDispatchersMetadata() *IMetadata {
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

func (m *IService[In, Out]) GetDispatcherMetadata(name string) *IMetadata {
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
