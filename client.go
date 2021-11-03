package elasticsteps

import (
	"context"
	"encoding/json"
)

// Client is an interface for interacting with Elasticsearch.
type Client interface {
	IndexGetter
	IndexCreator
	IndexDeleter
	DocumentIndexer
	DocumentFinder
}

// IndexGetter gets index.
type IndexGetter interface {
	GetIndex(ctx context.Context, index string) (json.RawMessage, error)
}

// IndexCreator creates indices.
type IndexCreator interface {
	CreateIndex(ctx context.Context, index string) error
	RecreateIndex(ctx context.Context, index string) error
}

// IndexDeleter deletes indices.
type IndexDeleter interface {
	DeleteIndex(ctx context.Context, indices ...string) error
}

// DocumentIndexer indexes documents.
type DocumentIndexer interface {
	IndexDocuments(ctx context.Context, index string, documents ...Document) error
}

// DocumentFinder gets documents.
type DocumentFinder interface {
	FindDocuments(ctx context.Context, index string, query *string) ([]json.RawMessage, error)
}
