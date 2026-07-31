package flow

import (
	"app/server/server/graph/model"
	"app/server/store"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

var (
	searchCounts = make(map[string]int)
	searchMutex  sync.Mutex
)

func getSearchCount(studentID string) int {
	searchMutex.Lock()
	defer searchMutex.Unlock()
	return searchCounts[studentID]
}

func incrementSearchCount(studentID string) {
	searchMutex.Lock()
	defer searchMutex.Unlock()
	searchCounts[studentID]++
}

// Input implements the struct as the param to pass in order to get the reply
type Input struct {
	context.Context
	*model.ChessStudentRequest
	*store.Store
}

// Flow main flow where the ai talks
func Flow(ctx context.Context, g *genkit.Genkit, input *model.ChessStudentRequest) (*model.OnChessCoachReply, error) {
	//
	//	queryEmbedding, err := embedText(ctx, *input.Query)
	//	if err != nil {
	//		return nil, fmt.Errorf("embedding query: %w", err)
	//	}

	// Use VectorSearch for first and third searches until milvus is properly filled
	//var chunks []*marlin.IScoredResult
	//searchCount := getSearchCount(*input.ID)
	//
	//	if searchCount == 1 || searchCount == 3 {
	//		chunks, err = store.VectorSearch(ctx, "ask-coach", "be-coach", queryEmbedding, 5)
	//	} else {
	//		chunks, err = store.HybridSearch(ctx, "ask-coach", "be-coach", *input.Query, queryEmbedding, 5)
	//	}
	//	incrementSearchCount(*input.ID)

	//if err != nil {
	//	return nil, fmt.Errorf("search: %w", err)
	//}

	prompt := prompt(input)

	resp, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}
	
	log.Println("replace: ", resp.Text())

	cleaned := cleanJSONResponse(resp.Text())

	var coach model.OnChessCoachReply
	if err := json.Unmarshal([]byte(cleaned), &coach); err != nil {
		return nil, fmt.Errorf("parse response: %w\nraw: %s", err, cleaned)
	}
	//	go func() {
	//		bgCtx := context.Background()
	//
	//		content := fmt.Sprintf("Q: %s\nA: %s", *input.Query, *coach.Information.Desc)
	//		embedding, err := embedText(bgCtx, content)
	//		if err != nil {
	//			return
	//		}
	//
	//		store.IngestIfNewWithEmbedding(bgCtx, content, *input.ID, "qa-history", map[string]any{
	//			"student": *input.Name,
	//			"query":   *input.Query,
	//			"type":    "qa",
	//		}, embedding)
	//	}()

	return &coach, nil
}
