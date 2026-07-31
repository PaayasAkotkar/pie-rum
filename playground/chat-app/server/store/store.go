package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	marlin "marlin/core"
	"strings"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// Store a store for pgvector
type Store struct {
	//db *pgxpool.Pool
	m   *marlin.Marlin
	cli *milvusclient.Client
}

func New(ctx context.Context) *Store {
	conf := &milvusclient.ClientConfig{
		Address:  "localhost:19530",
		Username: "minioadmin",
		Password: "minioadmin12",
	}
	cli, err := milvusclient.New(ctx, conf)
	if err != nil {
		panic(err)
	}

	m := marlin.New(cli)

	if !m.IsConnected(ctx) {
		cli.Close(ctx)
		return nil
	}

	log.Println("[marlin connected 🤗]")
	return &Store{
		m:   m,
		cli: cli,
	}
}

func (s *Store) HybridSearch(ctx context.Context, bucket, branch, queryText string, embd []float32, k int) ([]*marlin.IScoredResult, error) {
	return s.m.HybridSearch(ctx, bucket, branch, queryText, embd, k)
}

func (s *Store) VectorSearch(ctx context.Context, bucket, branch string, embd []float32, k int) ([]*marlin.IScoredResult, error) {
	return s.m.VectorSearch(ctx, bucket, branch, embd, k)
}

func (s *Store) CreateBucket(ctx context.Context, bucket string, dim int, p *marlin.IPush) error {
	return s.m.CreateBucket(ctx, bucket, dim, p)
}

func (s *Store) IngestIfNewWithEmbedding(ctx context.Context, bucket, branch, content string, metadata map[string]any, embedding []float32) error {
	return s.m.IngestIfNewWithEmbedding(ctx, bucket, branch, content, metadata, embedding)
}

func (s *Store) Close(ctx context.Context) error {
	if s.cli != nil {
		return s.cli.Close(ctx)
	}
	return nil
}

// ToContextString formats HybridSearch results into a single string
// suitable for dropping into an LLM prompt as retrieved context.
//
// Source is rendered as "<branch>/<name>". Content is the object's
// Text. Metadata always includes the reranked relevance Score, plus
// whatever keys Data decodes to if it happens to hold a JSON object -
// this package doesn't have a dedicated metadata field, so Data is a
// best-effort source for it. If Data isn't JSON (or is empty), only
// the score is shown.
func (s *Store) ToContextString(results []*marlin.IScoredResult) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("### RELEVANT KNOWLEDGE BASE ###\n")
	for i, r := range results {
		obj := r.Object.Object
		source := r.Object.Branch + "/" + obj.Name

		sb.WriteString(fmt.Sprintf("\n--- Source: %s [Result %d] ---\n", source, i+1))
		sb.WriteString(obj.Text)

		metaParts := []string{fmt.Sprintf("score:%.4f", r.Score)}
		if meta, ok := decodeMetadata(obj.Data); ok {
			for k, v := range meta {
				metaParts = append(metaParts, fmt.Sprintf("%s:%v", k, v))
			}
		}
		sb.WriteString("\nMetadata: ")
		sb.WriteString(strings.Join(metaParts, ", "))
		sb.WriteString("\n")
	}
	sb.WriteString("\n-------------------------------\n")
	return sb.String()
}

// decodeMetadata tries to interpret data as a JSON object. Returns
// ok=false if data is empty or isn't a JSON object - this is a guess
// about what's in the opaque payload field, not a guarantee.
func decodeMetadata(data []byte) (map[string]any, bool) {
	if len(data) == 0 {
		return nil, false
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, false
	}
	return m, true
}

// related to pg 😅 writeen by claude
//func New(databaseURL string) (*Store, error) {
//	config, err := pgxpool.ParseConfig(databaseURL)
//	if err != nil {
//		return nil, fmt.Errorf("parsing database url: %w", err)
//	}
//
//	pool, err := pgxpool.NewWithConfig(context.Background(), config)
//	if err != nil {
//		return nil, fmt.Errorf("connecting to postgres: %w", err)
//	}
//	return &Store{db: pool}, nil
//}
//
//func (s *Store) InitSchema(ctx context.Context) error {
//	if _, err := s.db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`); err != nil {
//		return fmt.Errorf("creating vector extension: %w", err)
//	}
//
//	var dim int
//	err := s.db.QueryRow(ctx, `
//		SELECT atttypmod
//		FROM pg_attribute
//		WHERE attrelid = 'chunks'::regclass AND attname = 'embedding'
//	`).Scan(&dim)
//
//	if err == nil && dim != 768 {
//		log.Printf("Dimension mismatch detected (%d != 768). Recreating chunks table...", dim)
//		s.db.Exec(ctx, "DROP TABLE IF EXISTS chunks")
//	}
//
//	_, err = s.db.Exec(ctx, `
//		CREATE TABLE IF NOT EXISTS chunks (
//			id        SERIAL PRIMARY KEY,
//			content   TEXT NOT NULL UNIQUE,
//			source    TEXT NOT NULL DEFAULT '',
//			metadata  JSONB DEFAULT '{}',
//			embedding vector(768)
//		)`)
//	if err != nil {
//		return fmt.Errorf("creating chunks table: %w", err)
//	}
//
//	_, err = s.db.Exec(ctx, `
//		CREATE INDEX IF NOT EXISTS chunks_embedding_idx
//		ON chunks USING hnsw (embedding vector_cosine_ops)`)
//	if err != nil {
//		return fmt.Errorf("creating vector index: %w", err)
//	}
//
//	_, err = s.db.Exec(ctx, `
//		CREATE INDEX IF NOT EXISTS chunks_fts_idx
//		ON chunks USING GIN (to_tsvector('english', content))`)
//	if err != nil {
//		return fmt.Errorf("creating fts index: %w", err)
//	}
//
//	return nil
//}
//
//func (s *Store) Ingest(ctx context.Context, content, source string, metadata map[string]any, embedding []float32) error {
//	if metadata == nil {
//		metadata = make(map[string]any)
//	}
//
//	var vec *pgvector.Vector
//	if embedding != nil {
//		v := pgvector.NewVector(embedding)
//		vec = &v
//	}
//
//	_, err := s.db.Exec(ctx, `
//		INSERT INTO chunks (content, source, metadata, embedding)
//		VALUES ($1, $2, $3, $4)
//		ON CONFLICT (content) DO NOTHING`,
//		content, source, metadata, vec)
//
//	return err
//}
//
//func (s *Store) IngestIfNewWithEmbedding(ctx context.Context, content, source string, metadata map[string]any, embedding []float32) error {
//	return s.Ingest(ctx, content, source, metadata, embedding)
//}
//
//func (s *Store) SearchByID(ctx context.Context, id int, topK int) ([]Chunk, error) {
//	var embedding pgvector.Vector
//	err := s.db.QueryRow(ctx, `SELECT embedding FROM chunks WHERE id = $1`, id).Scan(&embedding)
//	if err != nil {
//		return nil, fmt.Errorf("fetching embedding for id %d: %w", id, err)
//	}
//	return s.vectorSearch(ctx, embedding.Slice(), topK)
//}
//
//func (s *Store) HybridSearch(ctx context.Context, queryEmbedding []float32, queryText string, topK int) ([]Chunk, error) {
//	var vectorResults []Chunk
//	var err error
//
//	if queryEmbedding != nil {
//		vectorResults, err = s.vectorSearch(ctx, queryEmbedding, topK*2)
//		if err != nil {
//			log.Printf("Vector search error: %v", err)
//		}
//	}
//
//	keywordResults, err := s.keywordSearch(ctx, queryText, topK*2)
//	if err != nil {
//		log.Printf("Keyword search error: %v", err)
//	}
//
//	// Merge using Reciprocal Rank Fusion (RRF)
//	return mergeRRF(vectorResults, keywordResults, topK), nil
//}
//
//func (s *Store) vectorSearch(ctx context.Context, embedding []float32, limit int) ([]Chunk, error) {
//	rows, err := s.db.Query(ctx, `
//		SELECT id, content, source, metadata,
//		       1 - (embedding <=> $1) AS score
//		FROM chunks
//		ORDER BY embedding <=> $1
//		LIMIT $2`, pgvector.NewVector(embedding), limit)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var results []Chunk
//	for rows.Next() {
//		var c Chunk
//		if err := rows.Scan(&c.ID, &c.Content, &c.Source, &c.Metadata, &c.Score); err != nil {
//			return nil, err
//		}
//		results = append(results, c)
//	}
//	return results, rows.Err()
//}
//
//func (s *Store) keywordSearch(ctx context.Context, query string, limit int) ([]Chunk, error) {
//	if strings.TrimSpace(query) == "" {
//		return nil, nil
//	}
//
//	// ai generated
//	rows, err := s.db.Query(ctx, `
//		SELECT id, content, source, metadata,
//		       ts_rank(to_tsvector('english', content), plainto_tsquery('english', $1)) AS score
//		FROM chunks
//		WHERE to_tsvector('english', content) @@ plainto_tsquery('english', $1)
//		ORDER BY score DESC
//		LIMIT $2`, query, limit)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var results []Chunk
//	for rows.Next() {
//		var c Chunk
//		if err := rows.Scan(&c.ID, &c.Content, &c.Source, &c.Metadata, &c.Score); err != nil {
//			return nil, err
//		}
//		results = append(results, c)
//	}
//	return results, rows.Err()
//}
//
//func (s *Store) GetCount(ctx context.Context) (int, error) {
//	var count int
//	err := s.db.QueryRow(ctx, "SELECT count(*) FROM chunks").Scan(&count)
//	return count, err
//}
//
//func (s *Store) Clear(ctx context.Context) error {
//	_, err := s.db.Exec(ctx, "TRUNCATE TABLE chunks")
//	return err
//}
//
//func (s *Store) ToContextString(chunks []Chunk) string {
//	if len(chunks) == 0 {
//		return ""
//	}
//
//	var sb strings.Builder
//	sb.WriteString("### RELEVANT KNOWLEDGE BASE ###\n")
//	for i, c := range chunks {
//		sb.WriteString(fmt.Sprintf("\n--- Source: %s [Result %d] ---\n", c.Source, i+1))
//		sb.WriteString(c.Content)
//		if len(c.Metadata) > 0 {
//			sb.WriteString("\nMetadata: ")
//			metaParts := []string{}
//			for k, v := range c.Metadata {
//				metaParts = append(metaParts, fmt.Sprintf("%s:%v", k, v))
//			}
//			sb.WriteString(strings.Join(metaParts, ", "))
//		}
//		sb.WriteString("\n")
//	}
//	sb.WriteString("\n-------------------------------\n")
//	return sb.String()
//}
//
//// ai generated
//func mergeRRF(vectorResults, keywordResults []Chunk, topK int) []Chunk {
//	const k = 60.0 // Constant used in RRF to prevent low ranks from dominating
//	scores := make(map[int]float64)
//	chunkMap := make(map[int]Chunk)
//
//	for rank, c := range vectorResults {
//		scores[c.ID] += 1.0 / (k + float64(rank+1))
//		chunkMap[c.ID] = c
//	}
//	for rank, c := range keywordResults {
//		scores[c.ID] += 1.0 / (k + float64(rank+1))
//		if _, ok := chunkMap[c.ID]; !ok {
//			chunkMap[c.ID] = c
//		}
//	}
//
//	type ranked struct {
//		id    int
//		score float64
//	}
//	var all []ranked
//	for id, score := range scores {
//		all = append(all, ranked{id, score})
//	}
//
//	sort.Slice(all, func(i, j int) bool {
//		return all[i].score > all[j].score
//	})
//
//	results := make([]Chunk, 0, topK)
//	for i, r := range all {
//		if i >= topK {
//			break
//		}
//		c := chunkMap[r.id]
//		c.Score = r.score
//		results = append(results, c)
//	}
//	return results
//}
// end
