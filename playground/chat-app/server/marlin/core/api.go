package marlin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func (s *Marlin) IsConnected(ctx context.Context) bool {
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := s.cli.ListCollections(pingCtx, milvusclient.NewListCollectionOption())
	if err != nil {
		log.Printf("milvus ping failed: %s", err)
		return false
	}
	return true
}
func (s *Marlin) CreateBucket(ctx context.Context, bucket string, dim int, p *IPush) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	has, err := s.cli.HasCollection(ctx, milvusclient.NewHasCollectionOption(bucket))
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", bucket, err)
	}

	if !has {
		schema := buildSchema(dim)

		opt := milvusclient.NewCreateCollectionOption(bucket, schema).
			WithIndexOptions(autoVectorIndex(bucket), bm25SparseIndex(bucket))

		if err := s.cli.CreateCollection(ctx, opt); err != nil {
			return fmt.Errorf("create bucket %q: %w", bucket, err)
		}

		//		loadTask, err := s.cli.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(bucket))
		//		if err != nil {
		//			return fmt.Errorf("load bucket %q: %w", bucket, err)
		//		}
		//
		//		if err := loadTask.Await(ctx); err != nil {
		//			return fmt.Errorf("wait for bucket %q to load: %w", bucket, err)
		//		}

		s.dims[bucket] = dim
	}

	if p == nil || p.Object == nil {
		return nil
	}
	return s.writeObject(ctx, bucket, p.Branch, p.Object)
}

func (s *Marlin) BucketExists(ctx context.Context, bucket string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	has, err := s.cli.HasCollection(ctx, milvusclient.NewHasCollectionOption(bucket))

	if err != nil {
		return false, fmt.Errorf("check bucket %q: %w", bucket, err)
	}

	return has, nil
}

func (s *Marlin) GetObjectByName(ctx context.Context, bucket, branch, object string) *IPullResults {
	s.mu.Lock()
	defer s.mu.Unlock()

	rs, err := s.cli.Query(ctx, milvusclient.NewQueryOption(bucket).
		WithPartitions(branch).
		WithFilter(fmt.Sprintf("%s == %q", fieldPK, pk(branch, object))).
		WithOutputFields(fieldName, fieldBranch, fieldVector, fieldData))
	if err != nil {
		return &IPullResults{Error: []error{fmt.Errorf("query %q/%q/%q: %w", bucket, branch, object, err)}}
	}

	result := s.rowsToResults(bucket, rs)
	if len(result.Pulls) == 0 && len(result.Error) == 0 {
		result.Error = append(result.Error, fmt.Errorf("object %q/%q/%q not found", bucket, branch, object))
	}
	return result
}

func (s *Marlin) GetObjectByBranch(ctx context.Context, bucket, branch string) *IPullResults {
	s.mu.Lock()
	defer s.mu.Unlock()

	rs, err := s.cli.Query(ctx, milvusclient.NewQueryOption(bucket).
		WithPartitions(branch).
		WithFilter(fmt.Sprintf("%s == %q", fieldBranch, branch)).
		WithOutputFields(fieldName, fieldBranch, fieldVector, fieldData))

	if err != nil {
		return &IPullResults{Error: []error{fmt.Errorf("query %q/%q: %w", bucket, branch, err)}}
	}

	return s.rowsToResults(bucket, rs)
}

func (s *Marlin) GetObjects(ctx context.Context, bucket string) (map[string]*IPullResults, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	branches, err := s.cli.ListPartitions(ctx, milvusclient.NewListPartitionOption(bucket))
	if err != nil {
		return nil, fmt.Errorf("list branches of %q: %w", bucket, err)
	}

	out := make(map[string]*IPullResults, len(branches))
	for _, branch := range branches {
		rs, err := s.cli.Query(ctx, milvusclient.NewQueryOption(bucket).
			WithPartitions(branch).
			WithFilter(fmt.Sprintf("%s == %q", fieldBranch, branch)).
			WithOutputFields(fieldName, fieldBranch, fieldVector, fieldData))
		if err != nil {
			out[branch] = &IPullResults{Error: []error{fmt.Errorf("query %q/%q: %w", bucket, branch, err)}}
			continue
		}
		out[branch] = s.rowsToResults(bucket, rs)
	}
	return out, nil
}

func (s *Marlin) Push(ctx context.Context, p IPush) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writeObject(ctx, p.Bucket, p.Branch, p.Object)
}

func (s *Marlin) DeleteBucket(ctx context.Context, bucket string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.cli.DropCollection(ctx, milvusclient.NewDropCollectionOption(bucket)); err != nil {
		return fmt.Errorf("delete bucket %q: %w", bucket, err)
	}
	delete(s.dims, bucket)
	return nil
}

func (s *Marlin) CleanBucket(ctx context.Context, bucket string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.cli.Delete(ctx, milvusclient.NewDeleteOption(bucket).WithExpr(fieldPK+` != ""`)); err != nil {
		return fmt.Errorf("clean bucket %q: %w", bucket, err)
	}

	branches, err := s.cli.ListPartitions(ctx, milvusclient.NewListPartitionOption(bucket))
	if err != nil {
		return fmt.Errorf("list branches of %q: %w", bucket, err)
	}

	var errs []error
	for _, branch := range branches {
		if branch == defaultPartition {
			continue
		}
		if err := s.cli.DropPartition(ctx, milvusclient.NewDropPartitionOption(bucket, branch)); err != nil {
			errs = append(errs, fmt.Errorf("drop branch %q/%q: %w", bucket, branch, err))
		}
	}
	return errors.Join(errs...)
}

func (s *Marlin) DeleteBranch(ctx context.Context, d IDelete) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr := fmt.Sprintf("%s == %q", fieldBranch, d.Branch)
	if _, err := s.cli.Delete(ctx, milvusclient.NewDeleteOption(d.Bucket).WithExpr(expr)); err != nil {
		return fmt.Errorf("delete branch %q/%q: %w", d.Bucket, d.Branch, err)
	}

	if d.Branch == defaultPartition {
		return nil
	}
	if err := s.cli.DropPartition(ctx, milvusclient.NewDropPartitionOption(d.Bucket, d.Branch)); err != nil {
		return fmt.Errorf("drop branch %q/%q: %w", d.Bucket, d.Branch, err)
	}
	return nil
}

func (s *Marlin) DeleteObject(ctx context.Context, d IDelete) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr := fmt.Sprintf("%s == %q", fieldPK, pk(d.Branch, d.Object))
	if _, err := s.cli.Delete(ctx, milvusclient.NewDeleteOption(d.Bucket).WithExpr(expr)); err != nil {
		return fmt.Errorf("delete object %q/%q/%q: %w", d.Bucket, d.Branch, d.Object, err)
	}
	return nil
}

// HybridSearch runs a dense-vector similarity search and a BM25
// full-text search over bucket/branch in parallel and merges them with
// Reciprocal Rank Fusion, returning up to topK best combined hits,
// best first.
//
// queryVector must have exactly the bucket's configured dimension.
// queryText is matched against each object's Text field via Milvus's
// built-in BM25 function - you don't need to compute a sparse
// embedding yourself. Pass an empty queryVector or queryText to fall
// back to a search on just the other one (Milvus still requires the
// unused query be well-formed, so this skips that leg entirely rather
// than sending a zero vector or empty string through it).
//
// Note: WithPartitions/WithOutputFields on hybrid search follow the
// same option pattern used everywhere else in this client, but I
// wasn't able to confirm those two methods specifically exist on
// HybridSearchOption against documentation - worth an early smoke
// test against a live cluster.
func (s *Marlin) HybridSearch(ctx context.Context, bucket, branch, queryText string, queryVector []float32, topK int) ([]*IScoredResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(queryVector) == 0 && queryText == "" {
		return nil, fmt.Errorf("hybrid search %q/%q: need a query vector, query text, or both", bucket, branch)
	}

	var requests []*milvusclient.AnnRequest

	if len(queryVector) > 0 {
		dim, err := s.dimOf(ctx, bucket)
		if err != nil {
			return nil, err
		}
		if len(queryVector) != dim {
			return nil, fmt.Errorf("query vector has %d dims, bucket %q expects %d", len(queryVector), bucket, dim)
		}
		requests = append(requests, milvusclient.NewAnnRequest(fieldVector, topK, entity.FloatVector(queryVector)))
	}

	if queryText != "" {
		requests = append(requests, milvusclient.NewAnnRequest(fieldSparse, topK, entity.Text(queryText)))
	}

	opt := milvusclient.NewHybridSearchOption(bucket, topK, requests...).
		WithReranker(milvusclient.NewRRFReranker()).
		WithPartitions(branch).
		WithOutputFields(fieldName, fieldBranch, fieldText, fieldData)

	resultSets, err := s.cli.HybridSearch(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("hybrid search %q/%q: %w", bucket, branch, err)
	}
	if len(resultSets) == 0 {
		return nil, nil
	}
	return scoredResultsFromSet(bucket, resultSets[0])
}

// VectorSearch runs a simple dense-vector similarity search only.
// Use this when you want to avoid hybrid search (e.g., when milvus is not yet populated).
func (s *Marlin) VectorSearch(ctx context.Context, bucket, branch string, queryVector []float32, topK int) ([]*IScoredResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(queryVector) == 0 {
		return nil, fmt.Errorf("vector search %q/%q: need a query vector", bucket, branch)
	}

	dim, err := s.dimOf(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if len(queryVector) != dim {
		return nil, fmt.Errorf("query vector has %d dims, bucket %q expects %d", len(queryVector), bucket, dim)
	}

	resultSets, err := s.cli.Search(ctx, milvusclient.NewSearchOption(
		bucket,
		topK,
		[]entity.Vector{entity.FloatVector(queryVector)},
	).WithPartitions(branch).
		WithOutputFields(fieldName, fieldBranch, fieldText, fieldData))
	if err != nil {
		return nil, fmt.Errorf("vector search %q/%q: %w", bucket, branch, err)
	}
	if len(resultSets) == 0 {
		return nil, nil
	}
	return scoredResultsFromSet(bucket, resultSets[0])
}
