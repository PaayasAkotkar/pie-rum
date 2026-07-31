package pierum

//
//func (m *IMetadata) Rebuild(size int) {
//	m.Active = make([]MetadataInfo, 0, size)
//	m.InActive = make([]MetadataInfo, 0, size)
//	m.Swapped = make([]MetadataInfo, 0, size)
//	m.Rankings = make([]MetadataRankingInfo, 0, size)
//	m.Len = 0
//}
//
//func NewMetadata(size int) *IMetadata {
//	return &IMetadata{
//		Active:   make([]MetadataInfo, 0, size),
//		InActive: make([]MetadataInfo, 0, size),
//		Swapped:  make([]MetadataInfo, 0, size),
//		Rankings: make([]MetadataRankingInfo, 0, size),
//	}
//}
//
//func (m *IMetadata) JSON() []byte {
//	bytes, err := json.Marshal(m)
//	if err != nil {
//		log.Println(err)
//		return nil
//	}
//	return bytes
//}
//
//func (m *IMetadata) SaveLen() {
//	m.Len = int64(len(m.Rankings))
//}
//
//func (m *IMetadata) AddActive(info MetadataInfo) {
//	m.Active = append(m.Active, info)
//}
//
//func (m *IMetadata) AddInActive(info MetadataInfo) {
//	m.InActive = append(m.InActive, info)
//}
//
//func (m *IMetadata) AddSwapped(info MetadataInfo) {
//	m.Swapped = append(m.Swapped, info)
//}
//
//func (m *IMetadata) AddRanking(info MetadataRankingInfo) {
//	m.Rankings = append(m.Rankings, info)
//}
//
//func (m *MetadataInfo) AddUsageLen() {
//	m.UsageLen++
//}
//func (m *MetadataInfo) UpdateLastUsed(t time.Time) {
//	m.LastUpdated = common.FormatDateForClient(t)
//}
//func (m *MetadataInfo) UpdateName(name string) {
//	m.Name = name
//}
//
//func (r *MetadataRankingInfo) AddUsageLen() {
//	r.UsageLen++
//}
//func (r *MetadataRankingInfo) UpdateRank(rank int64) {
//	r.Rank = rank
//}
//func (r *MetadataRankingInfo) UpdateName(name string) {
//	r.Name = name
//}
