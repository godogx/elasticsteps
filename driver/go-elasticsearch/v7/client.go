package elasticsearch7

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bool64/ctxd"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"

	"github.com/godogx/elasticsteps"
)

var _ elasticsteps.Client = (*Client)(nil)

// Client is a wrapper around elasticsearch7.Client.
type Client struct {
	es *es7.Client
}

// GetIndex satisfies elasticsteps.Client.
func (c *Client) GetIndex(ctx context.Context, index string) (json.RawMessage, error) {
	get := c.es.Indices.Get

	_, err := refineResp(get([]string{index}, get.WithContext(ctx)))
	if err != nil {
		if err.code == http.StatusNotFound {
			return nil, elasticsteps.ErrIndexNotFound
		}

		return nil, err
	}

	return nil, nil
}

// CreateIndex satisfies elasticsteps.Client.
func (c *Client) CreateIndex(ctx context.Context, index string) error {
	create := c.es.Indices.Create

	_, err := refineResp(create(index, create.WithContext(ctx)))
	if err != nil {
		return ctxd.WrapError(ctx, err, "could not create index", "index", index)
	}

	return nil
}

// RecreateIndex satisfies elasticsteps.Client.
func (c *Client) RecreateIndex(ctx context.Context, index string) error {
	if err := c.DeleteIndex(ctx, index); err != nil {
		return err
	}

	return c.CreateIndex(ctx, index)
}

// DeleteIndex satisfies elasticsteps.Client.
func (c *Client) DeleteIndex(ctx context.Context, indices ...string) error {
	del := c.es.Indices.Delete

	_, err := refineResp(del(indices, del.WithContext(ctx)))
	if err != nil {
		if err.code == http.StatusNotFound {
			return nil
		}

		return ctxd.WrapError(ctx, err, "could not delete indices", "indices", indices)
	}

	return nil
}

// IndexDocuments satisfies elasticsteps.Client.
func (c *Client) IndexDocuments(ctx context.Context, index string, docs ...elasticsteps.Document) error {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:  c.es,
		Index:   index,
		Refresh: "true",
	})
	if err != nil {
		return ctxd.WrapError(ctx, err, "could not init bulk indexer", "index", index)
	}

	for _, doc := range docs {
		doc := doc

		err := indexer.Add(ctx, esutil.BulkIndexerItem{
			Index:      index,
			Action:     "index",
			DocumentID: doc.ID,
			Body:       bytes.NewReader(doc.Source),
		})
		if err != nil {
			return ctxd.WrapError(ctx, err, "could not add doc to bulk indexer",
				"index", index, "doc", doc,
			)
		}
	}

	if err := indexer.Close(ctx); err != nil {
		return ctxd.WrapError(ctx, err, "could not close bulk indexer",
			"index", index,
		)
	}

	stats := indexer.Stats()

	if stats.NumFailed > 0 {
		return ctxd.NewError(ctx, "could not index all documents",
			"num_docs", stats.NumRequests,
			"num_failure", stats.NumFailed,
		)
	}

	return nil
}

// FindDocuments satisfies elasticsteps.Client.
func (c *Client) FindDocuments(ctx context.Context, index string, query *string) ([]json.RawMessage, error) {
	search := c.es.Search

	var body string

	if query != nil && len(*query) > 0 {
		body = *query
	} else {
		body = `{"query": {"match_all":{}}}`
	}

	resp, err := refineResp(search(
		search.WithContext(ctx),
		search.WithIndex(index),
		search.WithBody(strings.NewReader(body)),
	))
	if err != nil {
		return nil, ctxd.WrapError(ctx, err, "could not get all documents", "index", index)
	}

	var result elasticsteps.SearchResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, ctxd.WrapError(ctx, err, "could not unmarshal all documents", "index", index)
	}

	return result.Hits.Hits, nil
}

func wrapClient(client *es7.Client) *Client {
	return &Client{es: client}
}

func refineResp(resp *esapi.Response, err error) (*esapi.Response, *err) {
	if err != nil {
		return nil, newError(codeUnknown, err.Error())
	}

	if resp.IsError() {
		return nil, newError(resp.StatusCode, resp.String())
	}

	return resp, nil
}
