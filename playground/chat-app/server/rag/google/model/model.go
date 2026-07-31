package googlemodel

import (
	"context"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

const (
	//apiKey      = "your-api-key"
	//geminiFlash = "gemini-2.5-flash"
	apiKey      = "your-api-key"
	geminiFlash = "googleai/gemini-2.0-flash" // ref: https://genkit.dev/
)

func GK(ctx context.Context) *genkit.Genkit {
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: apiKey}),
		genkit.WithDefaultModel(geminiFlash),
	)
	return g
}
