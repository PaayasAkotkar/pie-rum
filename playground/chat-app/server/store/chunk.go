package store

type Chunk struct {
	ID        int
	Content   string
	Source    string
	Metadata  map[string]any
	Score     float64
	Embedding []float32
}
