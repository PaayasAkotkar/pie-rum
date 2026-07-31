package marlin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// IngestIfNewWithEmbedding writes content into bucket/branch as a
// hybrid-searchable object, skipping the write if identical content
// was already ingested into that branch - mirroring the old Postgres
// store's `ON CONFLICT (content) DO NOTHING` dedup behavior.
//
// name is derived deterministically from content (a SHA-256 hash), so
// re-ingesting the same text always resolves to the same object rather
// than accumulating duplicates with random IDs.
//
// bucket must already exist (see CreateBucket) with a dimension
// matching len(embedding).
//
// Caveat: unlike Postgres's ON CONFLICT, this check-then-write isn't
// atomic - concurrent calls ingesting the same content at the same
// moment could both pass the existence check. If that's a real risk in
// your ingestion path, either serialize calls per bucket/branch, or
// skip the check entirely and just call Push (see below) - since name
// is a hash of content, a duplicate write is a harmless overwrite of
// itself, not a new row.
func (s *Marlin) IngestIfNewWithEmbedding(ctx context.Context, bucket, branch, content string, metadata map[string]any, embedding []float32) error {
	name := contentKey(content)

	exists, err := s.objectExists(ctx, bucket, branch, name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	return s.Push(ctx, IPush{
		Bucket: bucket,
		Branch: branch,
		Object: &IObject{
			Name:   name,
			Vector: embedding,
			Text:   content,
			Data:   data,
		},
	})
}

// Ingest is IngestIfNewWithEmbedding without the existence check -
// always upserts. Since name is a hash of content, ingesting the same
// content twice just overwrites the same row (e.g. with refreshed
// metadata), so this is the simpler choice unless you specifically
// need "first write wins" semantics.
func (s *Marlin) Ingest(ctx context.Context, bucket, branch, content string, metadata map[string]any, embedding []float32) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	return s.Push(ctx, IPush{
		Bucket: bucket,
		Branch: branch,
		Object: &IObject{
			Name:   contentKey(content),
			Vector: embedding,
			Text:   content,
			Data:   data,
		},
	})
}

// objectExists checks for a row by primary key without the "not
// found" error IngestIfNewWithEmbedding needs to treat as a normal,
// expected case.
func (s *Marlin) objectExists(ctx context.Context, bucket, branch, name string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rs, err := s.cli.Query(ctx, milvusclient.NewQueryOption(bucket).
		WithPartitions(branch).
		WithFilter(fmt.Sprintf("%s == %q", fieldPK, pk(branch, name))).
		WithOutputFields(fieldPK))
	if err != nil {
		return false, fmt.Errorf("check %q/%q/%q: %w", bucket, branch, name, err)
	}
	return rs.Len() > 0, nil
}

func contentKey(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
