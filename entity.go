package elasticsteps

import (
	"bytes"
	"encoding/json"
)

// Document represents an Elasticsearch doc.
// nolint: tagliatelle
type Document struct {
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

type document Document

// UnmarshalJSON compacts document body while marshaling.
func (d *Document) UnmarshalJSON(data []byte) error {
	var raw document

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	compactedSource := new(bytes.Buffer)

	if err := json.Compact(compactedSource, raw.Source); err != nil {
		return err
	}

	*d = Document(raw)
	d.Source = compactedSource.Bytes()

	return nil
}

// SearchResult represents the search result.
type SearchResult struct {
	Hits SearchResultHits
}

// SearchResultHits represents the hits.
// nolint: tagliatelle
type SearchResultHits struct {
	Total    SearchResultHitsTotal `json:"total"`
	MaxScore *float64              `json:"max_score"`
	Hits     SearchResultHitsHits  `json:"hits"`
}

// SearchResultHitsTotal represents the hits.total.
type SearchResultHitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// SearchResultHitsHits represents the hits.hits.
type SearchResultHitsHits []json.RawMessage

// UnmarshalJSON removes unwanted fields.
func (h *SearchResultHitsHits) UnmarshalJSON(data []byte) error {
	var raw []map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*h = make([]json.RawMessage, len(raw))

	for k, v := range raw {
		delete(v, "_index")
		(*h)[k], _ = json.Marshal(v) // nolint: errcheck
	}

	return nil
}
