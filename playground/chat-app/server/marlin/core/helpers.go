package marlin

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// Close releases the underlying Milvus connection.
func (s *Marlin) Close(ctx context.Context) error {
	return s.cli.Close(ctx)
}

// buildSchema defines the fixed collection layout used for every
// bucket: a composite string PK, filterable name/branch scalars, a
// fixed-dimension dense vector, a base64 blob payload, and a text
// field wired through a BM25 function into an auto-generated sparse
// vector for full-text search.
func buildSchema(dim int) *entity.Schema {
	bm25 := entity.NewFunction().
		WithName(bm25FuncName).
		WithInputFields(fieldText).
		WithOutputFields(fieldSparse).
		WithType(entity.FunctionTypeBM25)

	return entity.NewSchema().
		WithField(entity.NewField().WithName(fieldPK).WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(1024).WithIsPrimaryKey(true)).
		WithField(entity.NewField().WithName(fieldName).WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512)).
		WithField(entity.NewField().WithName(fieldBranch).WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512)).
		WithField(entity.NewField().WithName(fieldVector).WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(dim))).
		WithField(entity.NewField().WithName(fieldText).WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(maxTextLen).WithEnableAnalyzer(true)).
		WithField(entity.NewField().WithName(fieldSparse).WithDataType(entity.FieldTypeSparseVector)).
		WithField(entity.NewField().WithName(fieldData).WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(maxDataLen)).
		WithFunction(bm25)
}

// autoVectorIndex builds a default AUTOINDEX + cosine metric index for
// the dense vector field of bucket.
func autoVectorIndex(bucket string) milvusclient.CreateIndexOption {
	return milvusclient.NewCreateIndexOption(bucket, fieldVector, index.NewAutoIndex(entity.COSINE)).
		WithIndexName(fieldVector)
}

// bm25SparseIndex builds the SPARSE_INVERTED_INDEX/BM25 index required
// for full-text search on the auto-generated sparse field. 0.2 is
// Milvus's commonly-used default drop ratio for pruning low-weight
// terms at build time.
func bm25SparseIndex(bucket string) milvusclient.CreateIndexOption {
	return milvusclient.NewCreateIndexOption(bucket, fieldSparse, index.NewSparseInvertedIndex(entity.BM25, 0.2)).
		WithIndexName(fieldSparse)
}

// dimOf returns bucket's configured vector dimension, consulting the
// in-memory cache first and falling back to DescribeCollection (e.g.
// after a process restart, or for a bucket this instance didn't
// create). Caller must hold s.mu.
func (s *Marlin) dimOf(ctx context.Context, bucket string) (int, error) {
	if d, ok := s.dims[bucket]; ok {
		return d, nil
	}

	coll, err := s.cli.DescribeCollection(ctx, milvusclient.NewDescribeCollectionOption(bucket))
	if err != nil {
		return 0, fmt.Errorf("describe bucket %q: %w", bucket, err)
	}
	for _, f := range coll.Schema.Fields {
		if f.Name != fieldVector {
			continue
		}
		dim, err := f.GetDim()
		if err != nil {
			return 0, fmt.Errorf("read vector dimension for bucket %q: %w", bucket, err)
		}
		s.dims[bucket] = int(dim)
		return int(dim), nil
	}
	return 0, fmt.Errorf("bucket %q has no %q field - was it created by this package?", bucket, fieldVector)
}

// ensurePartition creates branch inside bucket if it doesn't already
// exist. Caller must hold s.mu.
func (s *Marlin) ensurePartition(ctx context.Context, bucket, branch string) error {
	has, err := s.cli.HasPartition(ctx, milvusclient.NewHasPartitionOption(bucket, branch))
	if err != nil {
		return fmt.Errorf("check branch %q/%q: %w", bucket, branch, err)
	}
	if has {
		return nil
	}
	if err := s.cli.CreatePartition(ctx, milvusclient.NewCreatePartitionOption(bucket, branch)); err != nil {
		return fmt.Errorf("create branch %q/%q: %w", bucket, branch, err)
	}
	return nil
}

// writeObject upserts obj into bucket/branch by its composite primary
// key. fieldSparse is deliberately never written here - Milvus
// generates it automatically from fieldText via the BM25 function, and
// supplying it yourself causes an "unexpected field" error on insert.
// Caller must hold s.mu.
func (s *Marlin) writeObject(ctx context.Context, bucket, branch string, obj *IObject) error {
	if obj == nil {
		return fmt.Errorf("nil object for %q/%q", bucket, branch)
	}

	dim, err := s.dimOf(ctx, bucket)
	if err != nil {
		return err
	}
	if len(obj.Vector) != dim {
		return fmt.Errorf("object %q: vector has %d dims, bucket %q expects %d", obj.Name, len(obj.Vector), bucket, dim)
	}

	if err := s.ensurePartition(ctx, bucket, branch); err != nil {
		return err
	}

	encoded, err := encodeData(obj.Data)
	if err != nil {
		return fmt.Errorf("object %q: %w", obj.Name, err)
	}

	cols := []column.Column{
		column.NewColumnVarChar(fieldPK, []string{pk(branch, obj.Name)}),
		column.NewColumnVarChar(fieldName, []string{obj.Name}),
		column.NewColumnVarChar(fieldBranch, []string{branch}),
		column.NewColumnFloatVector(fieldVector, dim, [][]float32{obj.Vector}),
		column.NewColumnVarChar(fieldText, []string{obj.Text}),
		column.NewColumnVarChar(fieldData, []string{encoded}),
	}

	opt := milvusclient.NewColumnBasedInsertOption(bucket, cols...).WithPartition(branch)
	if _, err := s.cli.Upsert(ctx, opt); err != nil {
		return fmt.Errorf("upsert %q/%q/%q: %w", bucket, branch, obj.Name, err)
	}
	return nil
}

// rowsToResults converts a query ResultSet into IPullResults, decoding
// the base64 payload column back into raw bytes for each row.
func (s *Marlin) rowsToResults(bucket string, rs milvusclient.ResultSet) *IPullResults {
	p := &IPullResults{Pulls: make([]*IPush, 0), Error: make([]error, 0)}

	names, err := varcharColumnData(rs.GetColumn(fieldName))
	if err != nil {
		p.Error = append(p.Error, fmt.Errorf("read %q column: %w", fieldName, err))
		return p
	}
	branches, err := varcharColumnData(rs.GetColumn(fieldBranch))
	if err != nil {
		p.Error = append(p.Error, fmt.Errorf("read %q column: %w", fieldBranch, err))
		return p
	}
	texts, err := varcharColumnData(rs.GetColumn(fieldText))
	if err != nil {
		p.Error = append(p.Error, fmt.Errorf("read %q column: %w", fieldText, err))
		return p
	}
	dataCol, err := varcharColumnData(rs.GetColumn(fieldData))
	if err != nil {
		p.Error = append(p.Error, fmt.Errorf("read %q column: %w", fieldData, err))
		return p
	}
	vecCol, err := floatVectorColumnData(rs.GetColumn(fieldVector))
	if err != nil {
		p.Error = append(p.Error, fmt.Errorf("read %q column: %w", fieldVector, err))
		return p
	}

	for i := range names {
		raw, err := decodeData(dataCol[i])
		if err != nil {
			p.Error = append(p.Error, fmt.Errorf("decode %q/%q/%q: %w", bucket, branches[i], names[i], err))
			continue
		}
		p.Pulls = append(p.Pulls, &IPush{
			Bucket: bucket,
			Branch: branches[i],
			Object: &IObject{Name: names[i], Vector: vecCol[i], Text: texts[i], Data: raw},
		})
	}
	return p
}

// scoredResultsFromSet converts one HybridSearch ResultSet into scored
// results. The dense vector isn't fetched back here (it's rarely
// needed after a search hit and doubles the response size); only Text
// and Data are populated on the returned objects.
func scoredResultsFromSet(bucket string, rs milvusclient.ResultSet) ([]*IScoredResult, error) {
	pks, err := varcharColumnData(rs.IDs)
	if err != nil {
		return nil, fmt.Errorf("read result ids: %w", err)
	}
	names, err := varcharColumnData(rs.GetColumn(fieldName))
	if err != nil {
		return nil, fmt.Errorf("read %q column: %w", fieldName, err)
	}
	branches, err := varcharColumnData(rs.GetColumn(fieldBranch))
	if err != nil {
		return nil, fmt.Errorf("read %q column: %w", fieldBranch, err)
	}
	texts, err := varcharColumnData(rs.GetColumn(fieldText))
	if err != nil {
		return nil, fmt.Errorf("read %q column: %w", fieldText, err)
	}
	dataCol, err := varcharColumnData(rs.GetColumn(fieldData))
	if err != nil {
		return nil, fmt.Errorf("read %q column: %w", fieldData, err)
	}

	out := make([]*IScoredResult, 0, len(pks))
	for i := range pks {
		raw, err := decodeData(dataCol[i])
		if err != nil {
			return nil, fmt.Errorf("decode %q/%q/%q: %w", bucket, branches[i], names[i], err)
		}
		out = append(out, &IScoredResult{
			Object: &IPush{
				Bucket: bucket,
				Branch: branches[i],
				Object: &IObject{Name: names[i], Text: texts[i], Data: raw},
			},
			Score: rs.Scores[i],
		})
	}
	return out, nil
}

func encodeData(b []byte) (string, error) {
	enc := base64.StdEncoding.EncodeToString(b)
	if len(enc) > maxDataLen {
		return "", fmt.Errorf("payload too large: %d encoded bytes exceeds %d-byte column limit", len(enc), maxDataLen)
	}
	return enc, nil
}

func decodeData(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func varcharColumnData(col column.Column) ([]string, error) {
	vc, ok := col.(*column.ColumnVarChar)
	if !ok {
		return nil, fmt.Errorf("unexpected column type: %T", col)
	}
	return vc.Data(), nil
}

func floatVectorColumnData(col column.Column) ([]entity.FloatVector, error) {
	vc, ok := col.(*column.ColumnFloatVector)
	if !ok {
		return nil, fmt.Errorf("unexpected column type: %T", col)
	}
	return vc.Data(), nil
}
