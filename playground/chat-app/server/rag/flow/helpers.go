package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// cleanJSONResponse from markdown code blocks
// claude said that it is necessary to do so for the ai to extract
func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	return text
}

// Request payload structure matching Hugging Face's OpenAI-compatible embeddings route
type EmbeddingRequest struct {
	Input interface{} `json:"input"` // Can be a string or a slice of strings
}

// EmbeddingResponse structures
type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// embedText written by the google's gemini
func embedText(ctx context.Context, text string) ([]float32, error) {
	//	g := googlemodel.GK(ctx)
	//
	//	embeddingModel := genkit.LookupEmbedder(g, "text-embedding-004")
	//
	//	em, err := genkit.Embed(ctx, g,
	//		ai.WithEmbedderName(embeddingModel.Name()),
	//		ai.WithDocs(ai.DocumentFromText(text, metadata)),
	//	)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//
	//	//res, err := em.Embed(ctx, &ai.EmbedRequest{
	//	//	Input: []*ai.Document{
	//	//		ai.DocumentFromText(text, nil),
	//	//	},
	//	//})
	//	if err != nil {
	//		return nil, fmt.Errorf("embed: %w", err)
	//	}
	//	if len(em.Embeddings) == 0 {
	//		return nil, fmt.Errorf("no embeddings returned")
	//	}
	//	return em.Embeddings[0].Embedding, nil

	hfToken := "your-api-key"

	url := "https://router.huggingface.co/hf-inference/models/google/embeddinggemma-300m/v1/embeddings"

	// Prepare input texts (EmbeddingGemma supports prefix formatting like task prompts)
	reqBody := EmbeddingRequest{
		Input: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+hfToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API Error (%d): %s\n", resp.StatusCode, string(bodyBytes))
		return nil, err
	}

	// Parse response
	var embResp EmbeddingResponse
	if err := json.Unmarshal(bodyBytes, &embResp); err != nil {
		panic(err)
	}

	// Output the resulting vector dimensions and values
	if len(embResp.Data) > 0 {
		vector := embResp.Data[0].Embedding
		fmt.Printf("Successfully generated vector! Dimension size: %d\n", len(vector))
		fmt.Printf("First 5 values: %v...\n", vector[:5])
		return embResp.Data[0].Embedding, nil
	}
	return nil, fmt.Errorf("not able to gen emb")
}
