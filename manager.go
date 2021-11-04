package elasticsteps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/swaggest/assertjson"
)

const defaultInstance = "_default"

// Manager manages the elasticsearch data.
type Manager struct {
	instances map[string]Client
	queries   map[string]map[string]*string
}

// nolint: ireturn
func (m *Manager) client(instance string) Client {
	return m.instances[instance]
}

func (m *Manager) registerPrerequisites(sc *godog.ScenarioContext) {
	sc.Step(`index "([^"]*)" is created in es "([^"]*)"$`, m.createIndex)
	sc.Step(`index "([^"]*)" is created$`, func(index string) error {
		return m.createIndex(index, defaultInstance)
	})

	sc.Step(`index "([^"]*)" is recreated in es "([^"]*)"$`, m.recreateIndex)
	sc.Step(`index "([^"]*)" is recreated$`, func(index string) error {
		return m.recreateIndex(index, defaultInstance)
	})

	sc.Step(`there is (?:an )?index "([^"]*)" in es "([^"]*)"$`, m.recreateIndex)
	sc.Step(`there is (?:an )?index "([^"]*)"$`, func(index string) error {
		return m.recreateIndex(index, defaultInstance)
	})

	sc.Step(`no index "([^"]*)" in es "([^"]*)"$`, m.deleteIndex)
	sc.Step(`no index "([^"]*)"$`, func(index string) error {
		return m.deleteIndex(index, defaultInstance)
	})

	sc.Step(`no rows in index "([^"]*)" of es "([^"]*)"$`, m.truncateIndex)
	sc.Step(`no rows in index "([^"]*)"$`, func(index string) error {
		return m.truncateIndex(index, defaultInstance)
	})

	sc.Step(`these docs are stored in index "([^"]*)" of es "([^"]*)"[:]?$`, m.indexDocs)
	sc.Step(`these docs are stored in index "([^"]*)"[:]?$`, func(index string, docs *godog.DocString) error {
		return m.indexDocs(index, defaultInstance, docs)
	})

	sc.Step(`docs (?:in|from) this file are stored in index "([^"]*)" of es "([^"]*)"[:]?$`, m.indexDocsFromFile)
	sc.Step(`docs (?:in|from) this file are stored in index "([^"]*)"[:]?$`, func(index string, body *godog.DocString) error {
		return m.indexDocsFromFile(index, defaultInstance, body)
	})
}

func (m *Manager) registerActions(sc *godog.ScenarioContext) {
	sc.Step(`I search in index "([^"]*)" of es "([^"]*)" with query[:]?$`, m.findDocuments)
	sc.Step(`I search in index "([^"]*)" with query[:]?$`, func(index string, query *godog.DocString) error {
		return m.findDocuments(index, defaultInstance, query)
	})
}

func (m *Manager) registerAssertions(sc *godog.ScenarioContext) {
	sc.Step(`index "([^"]*)" exists in es "([^"]*)"$`, m.assertIndexExists)
	sc.Step(`index "([^"]*)" exists$`, func(index string) error {
		return m.assertIndexExists(index, defaultInstance)
	})

	sc.Step(`index "([^"]*)" does not exist in es "([^"]*)"$`, m.assertIndexNotExists)
	sc.Step(`index "([^"]*)" does not exist$`, func(index string) error {
		return m.assertIndexNotExists(index, defaultInstance)
	})

	sc.Step(`no docs are available in index "([^"]*)" of es "([^"]*)"$`, m.assertNoDocs)
	sc.Step(`no docs are available in index "([^"]*)"$`, func(index string) error {
		return m.assertNoDocs(index, defaultInstance)
	})

	sc.Step(`only these docs are available in index "([^"]*)" of es "([^"]*)"[:]?$`, m.assertAllDocs)
	sc.Step(`only these docs are available in index "([^"]*)"[:]?$`, func(index string, docs *godog.DocString) error {
		return m.assertAllDocs(index, defaultInstance, docs)
	})

	sc.Step(`only docs (?:in|from) this file are available in index "([^"]*)" of es "([^"]*)"[:]?$`, m.assertAllDocsFromFile)
	sc.Step(`only docs (?:in|from) this file are available in index "([^"]*)"[:]?$`, func(index string, body *godog.DocString) error {
		return m.assertAllDocsFromFile(index, defaultInstance, body)
	})

	sc.Step(`these docs are found in index "([^"]*)" of es "([^"]*)"[:]?$`, m.assertFoundDocs)
	sc.Step(`these docs are found in index "([^"]*)"[:]?$`, func(index string, docs *godog.DocString) error {
		return m.assertFoundDocs(index, defaultInstance, docs)
	})

	sc.Step(`docs (?:in|from) this file are found in index "([^"]*)" of es "([^"]*)"[:]?$`, m.assertFoundDocsFromFile)
	sc.Step(`docs (?:in|from) this file are found in index "([^"]*)"[:]?$`, func(index string, body *godog.DocString) error {
		return m.assertFoundDocsFromFile(index, defaultInstance, body)
	})
}

// RegisterContext registers the manager to the test suite.
func (m *Manager) RegisterContext(sc *godog.ScenarioContext) {
	sc.Before(func(context.Context, *godog.Scenario) (context.Context, error) {
		m.queries = make(map[string]map[string]*string)

		return nil, nil
	})

	m.registerPrerequisites(sc)
	m.registerActions(sc)
	m.registerAssertions(sc)
}

func (m *Manager) createIndex(index, instance string) error {
	return m.client(instance).CreateIndex(context.Background(), index)
}

func (m *Manager) recreateIndex(index, instance string) error {
	return m.client(instance).RecreateIndex(context.Background(), index)
}

func (m *Manager) deleteIndex(index, instance string) error {
	return m.client(instance).DeleteIndex(context.Background(), index)
}

func (m *Manager) truncateIndex(index, instance string) error {
	return m.client(instance).RecreateIndex(context.Background(), index)
}

func (m *Manager) indexDocs(index, instance string, body *godog.DocString) error {
	var docs []Document

	if err := json.Unmarshal([]byte(body.Content), &docs); err != nil {
		return fmt.Errorf("could not read documents for indexing: %w", err)
	}

	return m.client(instance).IndexDocuments(context.Background(), index, docs...)
}

func (m *Manager) indexDocsFromFile(index, instance string, body *godog.DocString) error {
	content, err := os.ReadFile(body.Content)
	if err != nil {
		return fmt.Errorf("could not read docs from file %q: %w", body.Content, err)
	}

	return m.indexDocs(index, instance, &godog.DocString{Content: string(content)})
}

func (m *Manager) findDocuments(index, instance string, query *godog.DocString) error {
	if _, ok := m.queries[instance]; !ok {
		m.queries[instance] = make(map[string]*string)
	}

	m.queries[instance][index] = &query.Content

	return nil
}

func (m *Manager) assertIndexExists(index, instance string) error {
	_, err := m.client(instance).GetIndex(context.Background(), index)

	return err
}

func (m *Manager) assertIndexNotExists(index, instance string) error {
	_, err := m.client(instance).GetIndex(context.Background(), index)

	if errors.Is(err, ErrIndexNotFound) {
		return nil
	}

	return err
}

func (m *Manager) assertNoDocs(index, instance string) error {
	docs, err := m.client(instance).FindDocuments(context.Background(), index, nil)
	numDocs := len(docs)

	if numDocs > 0 {
		return fmt.Errorf("there are %d docs in index %q", numDocs, index) // nolint: goerr113
	}

	return err
}

func (m *Manager) assertAllDocs(index, instance string, body *godog.DocString) error {
	docs, err := m.client(instance).FindDocuments(context.Background(), index, nil)
	if err != nil {
		return err
	}

	actual, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	expected := []byte(body.Content)

	if err := assertjson.FailNotEqual(expected, actual); err != nil {
		return fmt.Errorf("failed to compare docs: %w", err)
	}

	return nil
}

func (m *Manager) assertAllDocsFromFile(index, instance string, body *godog.DocString) error {
	content, err := os.ReadFile(body.Content)
	if err != nil {
		return fmt.Errorf("could not read docs from file %q: %w", body.Content, err)
	}

	return m.assertAllDocs(index, instance, &godog.DocString{Content: string(content)})
}

func (m *Manager) assertFoundDocs(index, instance string, body *godog.DocString) error {
	var query *string

	if qs, ok := m.queries[instance]; ok {
		query = qs[index]
	}

	docs, err := m.client(instance).FindDocuments(context.Background(), index, query)
	if err != nil {
		return err
	}

	actual, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	expected := []byte(body.Content)

	if err := assertjson.FailNotEqual(expected, actual); err != nil {
		return fmt.Errorf("failed to compare docs: %w", err)
	}

	return nil
}

func (m *Manager) assertFoundDocsFromFile(index, instance string, body *godog.DocString) error {
	content, err := os.ReadFile(body.Content)
	if err != nil {
		return fmt.Errorf("could not read docs from file %q: %w", body.Content, err)
	}

	return m.assertFoundDocs(index, instance, &godog.DocString{Content: string(content)})
}

// ManagerOption sets up the manager.
type ManagerOption func(m *Manager)

// NewManager initiates a new data manager.
func NewManager(client Client, opts ...ManagerOption) *Manager {
	m := &Manager{
		instances: map[string]Client{
			defaultInstance: client,
		},
		queries: map[string]map[string]*string{},
	}

	for _, o := range opts {
		o(m)
	}

	return m
}

// WithInstance adds a new es instance.
func WithInstance(name string, client Client) ManagerOption {
	return func(m *Manager) {
		m.instances[name] = client
	}
}
