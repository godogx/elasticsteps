package elasticsteps

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type clientMocker func(tb testing.TB) *client

var _ Client = (*client)(nil)

// client is a elasticsteps.Client.
type client struct {
	mock.Mock
}

func (c *client) GetIndex(ctx context.Context, index string) (json.RawMessage, error) {
	results := c.Called(ctx, index)

	result := results.Get(0)
	err := results.Error(1)

	switch r := result.(type) {
	case nil:
		return nil, err

	case []byte:
		return r, err

	case string:
		return []byte(r), err
	}

	return result.(json.RawMessage), err
}

func (c *client) CreateIndex(ctx context.Context, index string, config *string) error {
	return c.Called(ctx, index, config).Error(0)
}

func (c *client) RecreateIndex(ctx context.Context, index string, config *string) error {
	return c.Called(ctx, index, config).Error(0)
}

func (c *client) DeleteIndex(ctx context.Context, indices ...string) error {
	i := 1
	args := make([]interface{}, i+len(indices))
	args[0] = ctx

	for _, idx := range indices {
		args[i] = idx
		i++
	}

	return c.Called(args...).Error(0)
}

func (c *client) IndexDocuments(ctx context.Context, index string, documents ...Document) error {
	i := 2
	args := make([]interface{}, i+len(documents))
	args[0] = ctx
	args[1] = index

	for _, doc := range documents {
		args[i] = doc
		i++
	}

	return c.Called(args...).Error(0)
}

func (c *client) FindDocuments(ctx context.Context, index string, query *string) ([]json.RawMessage, error) {
	results := c.Called(ctx, index, query)

	result := results.Get(0)
	err := results.Error(1)

	switch r := result.(type) {
	case nil:
		return nil, err

	case [][]byte:
		result := make([]json.RawMessage, len(r))

		for k, v := range r {
			result[k] = v
		}

		return result, err

	case []string:
		result := make([]json.RawMessage, len(r))

		for k, v := range r {
			result[k] = json.RawMessage(v)
		}

		return result, err

	case string:
		var result []json.RawMessage

		if err := json.Unmarshal([]byte(r), &result); err != nil {
			return nil, err
		}

		return result, nil
	}

	return result.([]json.RawMessage), err
}

func (c *client) DeleteAllDocuments(ctx context.Context, index string) error {
	return c.Called(ctx, index).Error(0)
}

// mockClient creates Client mock with cleanup to ensure all the expectations are met.
func mockClient(mocks ...func(c *client)) clientMocker {
	return func(tb testing.TB) *client {
		tb.Helper()

		c := &client{}

		for _, m := range mocks {
			m(c)
		}

		tb.Cleanup(func() {
			assert.True(tb, c.Mock.AssertExpectations(tb))
		})

		return c
	}
}
