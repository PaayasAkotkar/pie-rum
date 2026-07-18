package pierum

import (
	"encoding/json"
	"log"
	"pie-rum-sdk/common"
	"time"
)

type IMetadata struct {
	Active   []MetadataInfo        `json:"active"`
	InActive []MetadataInfo        `json:"inActive"`
	Swapped  []MetadataInfo        `json:"swapped"`
	Rankings []MetadataRankingInfo `json:"rankings"` // total number of metas + their sorted raking as per the usage
	Len      int64                 `json:"len"`      // total collected metadata len
}

func (m *IMetadata) Rebuild(size int) {
	m.Active = make([]MetadataInfo, 0, size)
	m.InActive = make([]MetadataInfo, 0, size)
	m.Swapped = make([]MetadataInfo, 0, size)
	m.Rankings = make([]MetadataRankingInfo, 0, size)
	m.Len = 0
}

func NewMetadata(size int) *IMetadata {
	return &IMetadata{
		Active:   make([]MetadataInfo, 0, size),
		InActive: make([]MetadataInfo, 0, size),
		Swapped:  make([]MetadataInfo, 0, size),
		Rankings: make([]MetadataRankingInfo, 0, size),
	}
}

func (m *IMetadata) JSON() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}

func (m *IMetadata) SaveLen() {
	m.Len = int64(len(m.Rankings))
}

func (m *IMetadata) AddActive(info MetadataInfo) {
	m.Active = append(m.Active, info)
}

func (m *IMetadata) AddInActive(info MetadataInfo) {
	m.InActive = append(m.InActive, info)
}

func (m *IMetadata) AddSwapped(info MetadataInfo) {
	m.Swapped = append(m.Swapped, info)
}

func (m *IMetadata) AddRanking(info MetadataRankingInfo) {
	m.Rankings = append(m.Rankings, info)
}

type MetadataInfo struct {
	Name        string `json:"name"`
	LastUpdated string `json:"lastUpdated"` // conv time
	UsageLen    int64  `json:"usageLen"`    // how many times been used
}

func (m *MetadataInfo) AddUsageLen() {
	m.UsageLen++
}
func (m *MetadataInfo) UpdateLastUsed(t time.Time) {
	m.LastUpdated = common.FormatDateForClient(t)
}
func (m *MetadataInfo) UpdateName(name string) {
	m.Name = name
}

type MetadataRankingInfo struct {
	Rank     int64  `json:"rank"`
	Name     string `json:"name"`
	UsageLen int64  `json:"usageLen"`
}

func (r *MetadataRankingInfo) AddUsageLen() {
	r.UsageLen++
}
func (r *MetadataRankingInfo) UpdateRank(rank int64) {
	r.Rank = rank
}
func (r *MetadataRankingInfo) UpdateName(name string) {
	r.Name = name
}
