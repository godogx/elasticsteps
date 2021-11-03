package bootstrap

import (
	es7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/godogx/elasticsteps"
	elasticsearch7 "github.com/godogx/elasticsteps/driver/go-elasticsearch/v7"
)

func newElasticsearch7(address string) (*elasticsteps.Manager, error) {
	cfg := es7.Config{
		Addresses: []string{address},
	}

	es, err := es7.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	m := elasticsearch7.NewManager(es,
		elasticsearch7.WithInstance(esExtra, es),
	)

	return m, nil
}
