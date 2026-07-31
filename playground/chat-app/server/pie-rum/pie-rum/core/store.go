// Package pierum ...
// flow -> implement the serach-engine easily pass the profile name -> get the bucket -> run the serach -> pass the res
package pierum

import (
	"context"
	"fmt"
	"pie-rum-sdk/stack"
	"sync"
)

// IStore the main thing you use
type IStore[In, Out any] struct {
	Registry map[string]*IProfile[In, Out] `json:"profile"`
	// search  *searchmanager.SearchManager
	// metrics        map[string]*IMetric
	result         []*IDispatchResult   `json:"-"`
	tags           *stack.Stack[string] `json:"-"` // only returns the active profile lists
	stack          *stack.Stack[string] `json:"-"`
	defaultProfile string               `json:"-"`
	//slate          *ISlate              `json:"-"`
	ctx context.Context
	mu  sync.Mutex
}

func NewStore[In, Out any](ctx context.Context) *IStore[In, Out] {
	return &IStore[In, Out]{
		tags:     stack.NewStack[string](),
		Registry: make(map[string]*IProfile[In, Out]),
		// metrics:  make(map[string]*IMetric),
		stack: stack.NewStack[string](),
		//slate: NewSlate(),
		// search:  search,
		ctx: ctx,
	}
}

func (r *IStore[In, Out]) nextKey() string {
	for _, ra := range r.stack.Max() {
		if rx, ok := r.Registry[ra]; ok {
			if rx.Config.getActivate() {
				return ra
			}
		}
	}
	return r.defaultProfile
}
func (s *IStore[In, Out]) SetDefaultProfile(name string) {
	s.defaultProfile = name
}

//func (k *IStore[In, Out]) AddProfileUsage(p string) {
//	k.mu.Lock()
//	defer k.mu.Unlock()
//	k.slate.usage[p]++
//}

func (s *IStore[In, Out]) ReplaceProfile(name string, profile *IProfile[In, Out]) {
	if _, ok := s.Registry[name]; ok {
		s.Registry[name] = profile
	}
}

// documentation
//
//func (m *IStore[In, Out]) GetProfilesMetadata() *IMetadata {
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
//// end
//
//func (m *IStore[In, Out]) GetProfileMetadata(name string) *IMetadata {
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

func NewStoreWithMonitorTags[In, Out any](ctx context.Context, tags []string) *IStore[In, Out] {
	t := stack.NewStack[string]()
	for _, r := range tags {
		t.Push(r)
	}
	return &IStore[In, Out]{
		tags:     t,
		Registry: make(map[string]*IProfile[In, Out]),
		// metrics:  make(map[string]*IMetric),
		ctx: ctx,
	}
}

func (s *IStore[In, Out]) getLists() []string {
	l := make([]string, 0, len(s.Registry))
	for r := range s.Registry {
		l = append(l, r)
	}
	return l
}

func (r *IStore[In, Out]) release() {
	// r.metrics = nil
	r.result = nil
}

func (k *IStore[In, Out]) AddResult(res *IDispatchResult) {
	// log.Println("pushing results: ", res)
	k.result = append(k.result, res)
	// log.Println("new result: ", res.Metric)
}

// get funcs

func (r *IStore[In, Out]) SetProfile(profile map[string]*IProfile[In, Out]) {
	r.Registry = profile
}

func (r *IStore[In, Out]) AddProfile(name string, profile *IProfile[In, Out]) {
	r.Registry[name] = profile
}

// end

func (s *IStore[In, Out]) SetMonitorTags(tags []string) {
	if len(s.Registry) == 0 {
		for _, r := range tags {
			s.tags.Push(r)
		}
		return
	}

	for _, r := range tags {
		if s.IsProfileActive(r) {
			if !s.tags.Contains(r) {
				s.tags.Push(r)
			}
		}
	}
}

func (s *IStore[In, Out]) GetTags() []string {
	return s.tags.Max()
}

// func (s *IStore[In, Out]) GetMetrics(name string) *IMetric { return s.metrics[name] }
func (s *IStore[in, out]) handleProfileActivation(key string) error {
	if d, ok := s.Registry[key]; ok {
		d.Config.setActivate(true)
		//s.slate.RecordChange(key)
		return nil
	}
	return nil
}

func (s *IStore[In, Out]) SetConfig(key string, config *IConfig) error {
	if d, ok := s.Registry[key]; ok {
		d.Config = config
		return nil
	}
	return nil
}

func (s *IStore[in, out]) handleProfileDeactivation(key string) error {
	if d, ok := s.Registry[key]; ok {
		d.Config.setActivate(false)
		//s.slate.RecordChange(key)
		return nil
	}
	return fmt.Errorf("profile %v not found or inactive", key)
}

func (s *IStore[in, out]) handleProfileSwap(key1, key2 string) error {
	if err := swap(s.Registry, s.Registry, key1, key2); err != nil {
		return err
	}
	//s.slate.RecordChange(key1)
	//s.slate.RecordChange(key2)
	return nil
}

func (p *IStore[in, out]) IsProfileActive(key string) bool {
	return p.Registry[key].Config.getActivate()
}

func (p *IStore[in, out]) IsProfileSwap(key string) *IConfig {
	if seq, ok := p.Registry[key]; ok {
		return seq.Config
	}
	return nil
}

//
//func (k *IStore[In, Out]) UpdateProfileSlateChange(name string) {
//	k.slate.RecordChange(name)
//}
//func (k *IStore[In, Out]) UpdateProfileSlateUsage(name string) {
//	k.slate.RecordUsage(name)
//}
