package elasticsteps_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/godogx/elasticsteps"
)

func TestDocument_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	const validPayload = `{
	"_id": "41",
	"_source": {
		"handle": "item-41",
		"name": "Item 41",
		"locale": "en_US"
	}
}`

	const validCompactedSource = `{"handle":"item-41","name":"Item 41","locale":"en_US"}`

	testCases := []struct {
		scenario       string
		payload        string
		expectedResult elasticsteps.Document
		expectedError  string
	}{
		{
			scenario:      "invalid payload",
			payload:       `{`,
			expectedError: `unexpected end of JSON input`,
		},
		{
			scenario:      "invalid id",
			payload:       `{"_id": 42}`,
			expectedError: `json: cannot unmarshal number into Go struct field document._id of type string`,
		},
		{
			scenario:      "missing document source",
			payload:       `{"_id": "42"}`,
			expectedError: `unexpected end of JSON input`,
		},
		{
			scenario: "success",
			payload:  validPayload,
			expectedResult: elasticsteps.Document{
				ID:     "41",
				Source: []byte(validCompactedSource),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var result elasticsteps.Document

			err := json.Unmarshal([]byte(tc.payload), &result)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestSearchResultHitsHit_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	const validPayload = `[{
	"_index" : "elasticsearch7_index_5",
	"_type" : "_doc",
	"_id" : "41",
	"_score" : 1.0,
	"_source" : {
	  "handle" : "item-41",
	  "name" : "Item 41",
	  "locale" : "en_US"
	}
}]`

	const expectedPayload = `{"_id":"41","_score":1.0,"_source":{"handle":"item-41","name":"Item 41","locale":"en_US"},"_type":"_doc"}`

	testCases := []struct {
		scenario       string
		payload        string
		expectedResult elasticsteps.SearchResultHitsHits
		expectedError  string
	}{
		{
			scenario:      "invalid payload",
			payload:       `{`,
			expectedError: `unexpected end of JSON input`,
		},
		{
			scenario:      "invalid data",
			payload:       "42",
			expectedError: `json: cannot unmarshal number into Go value of type []map[string]json.RawMessage`,
		},
		{
			scenario:       "valid payload",
			payload:        validPayload,
			expectedResult: elasticsteps.SearchResultHitsHits{[]byte(expectedPayload)},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var result elasticsteps.SearchResultHitsHits

			err := json.Unmarshal([]byte(tc.payload), &result)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
