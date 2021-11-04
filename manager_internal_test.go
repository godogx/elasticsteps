package elasticsteps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type managerMocker func(t *testing.T) *Manager

const (
	index    = "test-index"
	instance = "_default"
)

func TestManager_createIndex(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("CreateIndex", mock.Anything, mock.Anything).
					Return(errors.New("create error"))
			}),
			expected: errors.New("create error"),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("CreateIndex", context.Background(), index).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).createIndex(index, instance))
		})
	}
}

func TestManager_recreateIndex(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("RecreateIndex", mock.Anything, mock.Anything).
					Return(errors.New("recreate error"))
			}),
			expected: errors.New("recreate error"),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("RecreateIndex", context.Background(), index).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).recreateIndex(index, instance))
		})
	}
}

func TestManager_deleteIndex(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("DeleteIndex", mock.Anything, mock.Anything).
					Return(errors.New("delete error"))
			}),
			expected: errors.New("delete error"),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("DeleteIndex", context.Background(), index).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).deleteIndex(index, instance))
		})
	}
}

func TestManager_truncateIndex(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("RecreateIndex", mock.Anything, mock.Anything).
					Return(errors.New("recreate error"))
			}),
			expected: errors.New("recreate error"),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("RecreateIndex", context.Background(), index).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).truncateIndex(index, instance))
		})
	}
}

func TestManager_indexDocs(t *testing.T) {
	t.Parallel()

	const validPayload = `[
	{
		"_id": "41",
		"_source": {
			"handle": "item-41",
			"name": "Item 41",
			"locale": "en_US"
		}
	},
	{
		"_id": "42",
		"_source": {
			"handle": "item-42",
			"name": "Item 42",
			"locale": "en_US"
		}
	}
]`

	testCases := []struct {
		scenario      string
		mock          managerMocker
		payload       string
		expectedError string
	}{
		{
			scenario:      "invalid payload",
			mock:          mockManager(),
			payload:       `{`,
			expectedError: `could not read documents for indexing: unexpected end of JSON input`,
		},
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("IndexDocuments", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("index error"))
			}),
			payload:       validPayload,
			expectedError: `index error`,
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("IndexDocuments", context.Background(), index,
					Document{
						ID:     "41",
						Source: json.RawMessage(`{"handle":"item-41","name":"Item 41","locale":"en_US"}`),
					},
					Document{
						ID:     "42",
						Source: json.RawMessage(`{"handle":"item-42","name":"Item 42","locale":"en_US"}`),
					},
				).
					Return(errors.New("index error"))
			}),
			payload:       validPayload,
			expectedError: `index error`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mock(t).indexDocs(index, instance, &godog.DocString{Content: tc.payload})

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestManager_indexDocsFromFile_NotFound(t *testing.T) {
	t.Parallel()

	m := mockManager()(t)
	err := m.indexDocsFromFile(index, instance, &godog.DocString{Content: "unknown"})

	expected := `could not read docs from file "unknown": open unknown: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestManager_assertIndexExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("GetIndex", mock.Anything, mock.Anything).
					Return(nil, errors.New("get error"))
			}),
			expected: errors.New("get error"),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("GetIndex", context.Background(), index).
					Return(nil, nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).assertIndexExists(index, instance))
		})
	}
}

func TestManager_assertIndexNotExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		mock     managerMocker
		expected error
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("GetIndex", mock.Anything, mock.Anything).
					Return(nil, errors.New("get error"))
			}),
			expected: errors.New("get error"),
		},
		{
			scenario: "success when index does not exist",
			mock: mockManager(func(c *client) {
				c.On("GetIndex", context.Background(), index).
					Return(nil, ErrIndexNotFound)
			}),
		},
		{
			scenario: "success",
			mock: mockManager(func(c *client) {
				c.On("GetIndex", context.Background(), index).
					Return(nil, nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.mock(t).assertIndexNotExists(index, instance))
		})
	}
}

func TestManager_assertNoDocs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mock          managerMocker
		expectedError string
	}{
		{
			scenario: "failure",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("get error"))
			}),
			expectedError: "get error",
		},
		{
			scenario: "has documents",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return([]json.RawMessage{nil}, nil)
			}),
			expectedError: `there are 1 docs in index "test-index"`,
		},
		{
			scenario: "no documents",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", context.Background(), index, (*string)(nil)).
					Return([]json.RawMessage{}, nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mock(t).assertNoDocs(index, instance)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestManager_assertAllDocs(t *testing.T) {
	t.Parallel()

	payload41 := json.RawMessage(`{"handle":"item-41","name":"Item 41","locale":"en_US"}`)
	payload42 := json.RawMessage(`{"handle":"item-42","name":"Item 42","locale":"en_US"}`)

	testCases := []struct {
		scenario       string
		mock           managerMocker
		expectedResult string
		expectedError  string
	}{
		{
			scenario: "fail to get documents",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("get error"))
			}),
			expectedError: `get error`,
		},
		{
			scenario: "invalid payload",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return([]json.RawMessage{[]byte("}")}, nil)
			}),
			expectedError: "json: error calling MarshalJSON for type json.RawMessage: invalid character '}' looking for beginning of value",
		},
		{
			scenario: "not equal",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", context.Background(), index, (*string)(nil)).
					Return([]json.RawMessage{payload41, payload42}, nil)
			}),
			expectedResult: fmt.Sprintf("[%s]", payload41),
			expectedError: `failed to compare docs: not equal:
 [
   {
     "handle": "item-41",
     "locale": "en_US",
     "name": "Item 41"
   }
+  {
+    "handle": "item-42",
+    "locale": "en_US",
+    "name": "Item 42"
+  }
 ]
`,
		},
		{
			scenario: "equal",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", context.Background(), index, (*string)(nil)).
					Return([]json.RawMessage{payload41, payload42}, nil)
			}),
			expectedResult: fmt.Sprintf("[%s,%s]", payload41, payload42),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mock(t).assertAllDocs(index, instance, &godog.DocString{Content: tc.expectedResult})

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestManager_assertAllDocsFromFile_NotFound(t *testing.T) {
	t.Parallel()

	m := mockManager()(t)
	err := m.assertAllDocsFromFile(index, instance, &godog.DocString{Content: "unknown"})

	expected := `could not read docs from file "unknown": open unknown: no such file or directory`

	assert.EqualError(t, err, expected)
}

func TestManager_assertFoundDocs(t *testing.T) {
	t.Parallel()

	payload41 := json.RawMessage(`{"handle":"item-41","name":"Item 41","locale":"en_US"}`)
	payload42 := json.RawMessage(`{"handle":"item-42","name":"Item 42","locale":"en_US"}`)

	testCases := []struct {
		scenario       string
		mock           managerMocker
		expectedResult string
		expectedError  string
	}{
		{
			scenario: "fail to get documents",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("get error"))
			}),
			expectedError: `get error`,
		},
		{
			scenario: "invalid payload",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", mock.Anything, mock.Anything, mock.Anything).
					Return([]json.RawMessage{[]byte("}")}, nil)
			}),
			expectedError: "json: error calling MarshalJSON for type json.RawMessage: invalid character '}' looking for beginning of value",
		},
		{
			scenario: "not equal",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", context.Background(), index, (*string)(nil)).
					Return([]json.RawMessage{payload41, payload42}, nil)
			}),
			expectedResult: fmt.Sprintf("[%s]", payload41),
			expectedError: `failed to compare docs: not equal:
 [
   {
     "handle": "item-41",
     "locale": "en_US",
     "name": "Item 41"
   }
+  {
+    "handle": "item-42",
+    "locale": "en_US",
+    "name": "Item 42"
+  }
 ]
`,
		},
		{
			scenario: "equal",
			mock: mockManager(func(c *client) {
				c.On("FindDocuments", context.Background(), index, (*string)(nil)).
					Return([]json.RawMessage{payload41, payload42}, nil)
			}),
			expectedResult: fmt.Sprintf("[%s,%s]", payload41, payload42),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mock(t).assertFoundDocs(index, instance, &godog.DocString{Content: tc.expectedResult})

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestManager_assertFoundDocs_WithQuery(t *testing.T) {
	t.Parallel()

	payload41 := json.RawMessage(`{"handle":"item-41","name":"Item 41","locale":"en_US"}`)
	payload42 := json.RawMessage(`{"handle":"item-42","name":"Item 42","locale":"en_US"}`)

	query := `{"query": {"match": {}}}`
	expected := fmt.Sprintf("[%s,%s]", payload41, payload42)

	m := mockManager(func(c *client) {
		c.On("FindDocuments", context.Background(), index, &query).
			Return([]json.RawMessage{payload41, payload42}, nil)
	})(t)

	err := m.findDocuments(index, instance, &godog.DocString{Content: query})
	assert.NoError(t, err)

	err = m.assertFoundDocs(index, instance, &godog.DocString{Content: expected})
	assert.NoError(t, err)
}

func TestManager_assertFoundDocsFromFile_NotFound(t *testing.T) {
	t.Parallel()

	m := mockManager()(t)
	err := m.assertFoundDocsFromFile(index, instance, &godog.DocString{Content: "unknown"})

	expected := `could not read docs from file "unknown": open unknown: no such file or directory`

	assert.EqualError(t, err, expected)
}

func mockManager(mocks ...func(c *client)) func(t *testing.T) *Manager {
	return func(t *testing.T) *Manager {
		t.Helper()

		return NewManager(mockClient(mocks...)(t))
	}
}
