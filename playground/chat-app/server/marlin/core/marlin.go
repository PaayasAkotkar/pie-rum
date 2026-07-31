// Package marlin mirrors the squid package's bucket/branch/object API,
// but backed by Milvus instead of Ceph RADOS:
//
//	Bucket  -> Milvus collection
//	Branch  -> Milvus partition inside that collection
//	Object  -> one row, holding a fixed-dimension vector, a base64 blob
//	           payload, and a text field usable for hybrid RAG search
//
// A few things don't map cleanly and are worth knowing before you use
// this:
//
//   - Milvus primary keys are unique per *collection*, not per
//     partition, so an object's real key is "<branch>/<name>", not just
//     "<name>". Two branches can each have their own "foo" object.
//   - Every bucket needs a fixed vector dimension, chosen when the
//     collection is created. All objects written to that bucket must
//     supply a vector of that exact length.
//   - Milvus has no "append" the way RADOS objects do. Push and
//     CreateBucket's object-write both upsert by primary key: writing
//     the same branch/name again replaces it.
//   - Data is base64-encoded into a VarChar column capped at 65535
//     characters, i.e. roughly 48KB of raw payload. This is not a
//     general-purpose object store; for large blobs, store them
//     elsewhere and keep only a reference here.
//   - Text is separate from Data: it's plain (unencoded) text, indexed
//     by Milvus's built-in BM25 function into an auto-generated sparse
//     vector, so HybridSearch can full-text-match on it alongside the
//     dense vector. Leave it empty for objects you don't need to
//     full-text search.
//   - Every collection has a built-in "_default" partition that Milvus
//     won't let you drop; it will keep showing up in branch listings.
//   - Writes are not guaranteed instantly visible to reads depending on
//     your consistency level; this package doesn't tune that for you.
//
// ~ written by claude ...
package marlin

import (
	"sync"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// Marlin wraps a Milvus client. mu serializes all access, and dims
// caches each bucket's configured vector dimension so repeated writes
// don't need a DescribeCollection round trip.
type Marlin struct {
	cli  *milvusclient.Client
	mu   sync.Mutex
	dims map[string]int
}

func New(cli *milvusclient.Client) *Marlin {
	return &Marlin{
		cli:  cli,
		dims: make(map[string]int),
	}
}

// IPush describes a write of a single object into a bucket/branch.
type IPush struct {
	Bucket string
	Branch string
	Object *IObject
}

// IObject is the payload written to, or read from, Milvus.
type IObject struct {
	Name   string
	Vector []float32 // must have exactly the bucket's configured dimension
	Text   string    // plain text, BM25-indexed for HybridSearch; may be empty
	Data   []byte    // arbitrary payload; stored base64-encoded, capped at maxDataLen encoded bytes
}

// IPullResults collects the objects (and any errors) returned by a read.
type IPullResults struct {
	Pulls []*IPush
	Error []error
}

// IScoredResult is one hit from a HybridSearch: an object plus its
// reranked relevance score (higher is more relevant).
type IScoredResult struct {
	Object *IPush
	Score  float32
}

// IDelete addresses a bucket, branch or object to be removed.
type IDelete struct {
	Bucket, Branch, Object string
}

// Schema field/function names used for every bucket.
const (
	fieldPK     = "pk" // "<branch>/<name>", primary key
	fieldName   = "name"
	fieldBranch = "branch"
	fieldVector = "vector" // dense embedding
	fieldText   = "text"   // plain text, input to the BM25 function
	fieldSparse = "sparse" // sparse vector auto-generated from fieldText by BM25 - never written directly
	fieldData   = "data"   // base64 blob payload

	bm25FuncName = "text_bm25_emb"

	// maxDataLen is the VarChar column limit for the base64-encoded
	// blob payload. Raw payloads must fit in roughly 3/4 of this after
	// encoding.
	maxDataLen = 65535

	// maxTextLen is the VarChar column limit for the plain-text field
	// used for BM25 full-text search.
	maxTextLen = 65535

	// defaultPartition is Milvus's built-in partition that always
	// exists and can't be dropped.
	defaultPartition = "_default"
)

// pk builds the composite primary key for an object, since Milvus PKs
// must be unique per collection rather than per partition.
func pk(branch, name string) string {
	return branch + "/" + name
}
