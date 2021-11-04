# Cucumber ElasticSearch steps for Golang

[![GitHub Releases](https://img.shields.io/github/v/release/godogx/elasticsteps)](https://github.com/godogx/elasticsteps/releases/latest)
[![Build Status](https://github.com/godogx/elasticsteps/actions/workflows/test.yaml/badge.svg)](https://github.com/godogx/elasticsteps/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/godogx/elasticsteps/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/godogx/elasticsteps)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/httpmock)](https://goreportcard.com/report/github.com/nhatthm/httpmock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/godogx/elasticsteps)

`elasticsteps` provides steps for [`cucumber/godog`](https://github.com/cucumber/godog) and makes it easy to run tests with ElasticSearch.

## Prerequisites

- `Go >= 1.17`

## Usage

### Setup

Initiate an `elasticsteps.Manager` and register it to the scenario using one of the supported drivers:

| Driver | Constructor |
| :--- | :---: |
| [`elastic/go-elasticsearch/v7`](https://github.com/elastic/go-elasticsearch) | [`driver/go-elasticsearch/v7`](https://github.com/godogx/elasticsteps/blob/master/driver/go-elasticsearch/v7/manager.go#L10) |
| [`olivere/elastic`](https://github.com/olivere/elastic) | ❌ |

```go
package mypackage

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cucumber/godog"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/require"

	elasticsearch7 "github.com/godogx/elasticsteps/driver/go-elasticsearch/v7"
)

func TestIntegration(t *testing.T) {
	out := bytes.NewBuffer(nil)

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"127.0.0.1:9200"},
	})
	require.NoError(t, err)

	// Create a new grpc client.
	manager := elasticsearch7.NewManager(es,
		// If you have another server, you could register it as well, for example:
		// elasticsearch7.WithInstance("another_instance", es),
	)

	suite := godog.TestSuite{
		Name:                 "Integration",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			// Register the client.
			manager.RegisterContext(ctx)
		},
		Options: &godog.Options{
			Strict:    true,
			Output:    out,
			Randomize: rand.Int63(),
		},
	}

	// Run the suite.
	if status := suite.Run(); status != 0 {
		t.Fatal(out.String())
	}
}
```

### Steps

#### Create a new index

Create a new index in the instance. If the index exists, the manager will throw an error.

- `index "([^"]*)" is created$`
- `index "([^"]*)" is created in es "([^"]*)"$` (if you want to create in the other instance)

For example:

```gherkin
Given index "products" is created
```

You could create a new index with a custom config by running:
 - `index "([^"]*)" is created with config[:]?$`
 - `index "([^"]*)" is created in es "([^"]*)" with config[:]?$`
 - `index "([^"]*)" is created with config from file[:]?$`
 - `index "([^"]*)" is created in es "([^"]*)" with config from file[:]?$`

For example:

```gherkin
Given index "products" is created with config:
"""
{
    "mappings": {
        "properties": {
            "size": {
                "type": "integer"
            },
            "name": {
                "type": "keyword"
            },
            "description": {
                "type": "text"
            }
        }
    }
}
"""
```

or

```gherkin
Given index "products" is created with config from file:
"""
../../resources/fixtures/mapping.json
"""
```

#### Recreate an index

Create a new index in the instance if it does not exist, otherwise the index will be deleted and recreated.

- `index "([^"]*)" is recreated$`
- `there is (?:an )?index "([^"]*)"$`
- `index "([^"]*)" is recreated in es "([^"]*)"$` (if you want to recreate in the other instance)
- `there is (?:an )?index "([^"]*)" in es "([^"]*)"$` (if you want to recreate in the other instance)

For example:

```gherkin
Given index "products" is recreated
```

or

```gherkin
Given there is index "products"
```

Same as [#Create a new index](#Create a new index), you could recreate an index with a custom config:
- Inline
  - `index "([^"]*)" is recreated with config[:]?$`
  - `there is (?:an )?index "([^"]*)" with config[:]?$`
- From a file
  - `index "([^"]*)" is recreated with config from file[:]?$`
  - `there is (?:an )?index "([^"]*)" with config from file[:]?$`
- Inline (for other instances)
  - `index "([^"]*)" is recreated in es "([^"]*)" with config[:]?$`
  - `there is (?:an )?index "([^"]*)" in es "([^"]*)" with config[:]?$`
- From a file (for other instances)
  - `index "([^"]*)" is recreated in es "([^"]*)" with config from file[:]?$`
  - `there is (?:an )?index "([^"]*)" in es "([^"]*)" with config from file[:]?$`

#### Delete an index

Delete an index in the instance. If the index does not exist, the manager will throw an error.

- `no index "([^"]*)"$`
- `no index "([^"]*)" in es "([^"]*)"$` (if you want to delete in the other instance)

For example:

```gherkin
Given no index "products"
```

#### Index Documents

- `these docs are stored in index "([^"]*)"[:]?$`
- `these docs are stored in index "([^"]*)" of es "([^"]*)"[:]?$` (if you want to index in the other instance)

For example:

```gherkin
Given these docs are stored in index "products":
"""
[
    {
        "_id": "41",
        "_source": {
            "handle": "item-41",
            "name": "Item 41",
            "locale": "en_US"
        }
    }
]
"""
```

You can also send the docs from a file by using:
- `docs (?:in|from) this file are stored in index "([^"]*)"[:]?$`
- `docs (?:in|from) this file are stored in index "([^"]*)" of es "([^"]*)"[:]?$` (if you want to index in the other instance)

For example:

```gherkin
Given docs in this file are stored in index "products":
"""
../../resources/fixtures/products.json
"""
```

#### Check whether an index exists

- `index "([^"]*)" exists$`
- `index "([^"]*)" exists in es "([^"]*)"$` (if you want to check the other instance)

For example:

```gherkin
Then index "products" exists
```

#### Check whether an index does not exist

- `index "([^"]*)" does not exist$`
- `index "([^"]*)" does not exist in es "([^"]*)"$` (if you want to check the other instance)

For example:

```gherkin
Then index "products" does not exist
```

#### Check there is no document in the index

- `no docs are available in index "([^"]*)"$`
- `no docs are available in index "([^"]*)" of es "([^"]*)"$` (if you want to check the other instance)

For example:

```gherkin
Then no docs are available in index "products"
```

#### Check whether index contains the exact documents

- `only these docs are available in index "([^"]*)"[:]?$`
- `only these docs are available in index "([^"]*)" of es "([^"]*)"[:]?$` (if you want to check the other instance)

For example:

```gherkin
Then only these docs are available in index "products":
"""
[
    {
        "_id": "41",
        "_source": {
            "handle": "item-41",
            "name": "Item 41",
            "locale": "en_US"
        },
        "_score": 1,
        "_type": "_doc"
    },
    {
        "_id": "42",
        "_source": {
            "handle": "item-42",
            "name": "Item 42",
            "locale": "en_US"
        },
        "_score": 1,
        "_type": "_doc"
    }
]
"""
```

You can also get the expected docs from a file by using:
- `only docs (?:in|from) this file are available in index "([^"]*)"[:]?$`
- `only docs (?:in|from) this file are available in index "([^"]*)" of es "([^"]*)"[:]?$` (if you want to index in the other instance)

For example:

```gherkin
Given only docs in this file are available in index "products":
"""
../../resources/fixtures/result.json
"""
```

#### Query documents

- First step: Setup the query <br/>
  `I search in index "([^"]*)" with query[:]?$` <br/>
  `I search in index "([^"]*)" of es "([^"]*)" with query[:]?$`
- Second step: Check the result <br/>
  `these docs are found in index "([^"]*)"[:]?$` <br/>
  `these docs are found in index "([^"]*)" of es "([^"]*)"[:]?$`

For example:

```gherkin
When I search in index "products" with query:
"""
{
    "query": {
        "match": {
            "locale": "en_US"
        }
    }
}
"""

Then these docs are found in index "products":
"""
[
    {
        "_id": "41",
        "_source": {
            "handle": "item-41",
            "name": "Item 41",
            "locale": "en_US"
        },
        "_score": "<ignore-diff>",
        "_type": "_doc"
    },
    {
        "_id": "42",
        "_source": {
            "handle": "item-42",
            "name": "Item 42",
            "locale": "en_US"
        },
        "_score": "<ignore-diff>",
        "_type": "_doc"
    }
]
"""
```

You can also get the expected docs from a file by using:
- `docs (?:in|from) this file are found in index "([^"]*)"[:]?$`
- `docs (?:in|from) this file are found in index "([^"]*)" of es "([^"]*)"[:]?$` (if you want to index in the other instance)

For example:

```gherkin
Given docs in this file are found in index "products":
"""
../../resources/fixtures/result.json
"""
```
