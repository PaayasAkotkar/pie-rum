package pierum

type ISequence[In any] struct {
	Profile string // to serach in profile
	// Service string // to manipulate Service
	// Rank  int64 // sequence number whether 1 ,.....
	Input *In
	// config *IConfig
}
