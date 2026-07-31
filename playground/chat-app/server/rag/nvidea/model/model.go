package nvideamodel

import (
	"context"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai"
)

// if you dont prefer the open-router simply just
// find the model that can be run using the open-ai
// and pass it here
const (
	//openRouterAPIKey = "your-api-key"
	//nvideaModel      = "nvidia/nemotron-nano-12b-v2-vl:free"
	openRouterAPIKey = "your-api-key"
	nvideaModel      = "openrouter/nvidia/nemotron-nano-12b-v2-vl:free"
	nvideaModel2     = "openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free"
	provider         = "openrouter"
)

func GK(ctx context.Context) *genkit.Genkit {

	oai := &compat_oai.OpenAICompatible{
		Provider: provider,
		BaseURL:  "https://openrouter.ai/api/v1",
		APIKey:   openRouterAPIKey,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(oai),
		genkit.WithDefaultModel(nvideaModel),
	)

	return g
}

func GKOmni(ctx context.Context) *genkit.Genkit {
	oai := &compat_oai.OpenAICompatible{
		Provider: provider,
		BaseURL:  "https://openrouter.ai/api/v1",
		APIKey:   openRouterAPIKey,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(oai),
		genkit.WithDefaultModel(nvideaModel2),
	)

	return g
}
