// Package pierum...
// schema.go implements the write to main entry where the metadata is exchange
package pierum

// required by plugin

type IPluginPackage struct {
	Result *IResults `json:"results"`
	Org    string    `json:"org"`
	Doc    string    `json:"doc"` // contains the info of the pie-rum core

	//Metaboard []IMetadataInfo `json:"info"` // on hold
}

// end
//
//type IMetadataInfo struct {
//	Profile   string         `json:"key"`
//	MetaBoard IMetadataBoard `json:"board"`
//}
//
//type IMetadataBoard struct {
//	Profile    []*IMetadata `json:"profile"`
//	Kit        []*IMetadata `json:"kit"`
//	Service    []*IMetadata `json:"service"`
//	Dispatcher []*IMetadata `json:"dispatcher"`
//	Events     []*IMetadata `json:"events"`
//}
//
//type IMetadata struct {
//	Active   []MetadataInfo        `json:"active"`
//	InActive []MetadataInfo        `json:"inActive"`
//	Swapped  []MetadataInfo        `json:"swapped"`
//	Rankings []MetadataRankingInfo `json:"rankings"` // total number of metas + their sorted raking as per the usage
//	Len      int64                 `json:"len"`      // total collected metadata len
//}
//
//type MetadataInfo struct {
//	Name        string `json:"name"`
//	LastUpdated string `json:"lastUpdated"` // conv time
//	UsageLen    int64  `json:"usageLen"`    // how many times been used
//}
//
//type MetadataRankingInfo struct {
//	Rank     int64  `json:"rank"`
//	Name     string `json:"name"`
//	UsageLen int64  `json:"usageLen"`
//}

// required by non-plguin

type IResults struct {
	IsReady bool               `json:"is_ready"`
	Resuts  []*IDispatchResult `json:"results"`
}

// end
