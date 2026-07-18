package pierum


// todo
// import (
// 	"context"
// 	"fmt"

// 	fastembed "github.com/anush008/fastembed-go"
// )

// // Embd wraps a fastembed engine for text embedding
// type Embd struct {
// 	engine *fastembed.FlagEmbedding
// 	dim    int
// }

// func NewFastEmbedder(cacheDir string) (*Embd, error) {
// 	show := true
// 	eng, err := fastembed.NewFlagEmbedding(&fastembed.InitOptions{
// 		Model:                fastembed.AllMiniLML6V2,
// 		MaxLength:            200,
// 		CacheDir:             cacheDir,
// 		ShowDownloadProgress: &show,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("fastembed init: %w", err)
// 	}
// 	return &Embd{engine: eng, dim: 384}, nil
// }

// func (f *Embd) Embed(ctx context.Context, text string) ([]float32, error) {
// 	tres, err := f.engine.Embed([]string{text}, 32)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(tres) == 0 {
// 		return nil, fmt.Errorf("fastembed returned empty result")
// 	}
// 	out := make([]float32, len(tres[0]))
// 	for i, v := range tres[0] {
// 		out[i] = float32(v)
// 	}
// 	return out, nil
// }

// func (f *Embd) Dim() int { return f.dim }
