package elasticsearch7

import (
	es7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/godogx/elasticsteps"
)

// NewManager initiates a new data manager.
func NewManager(client *es7.Client, opts ...elasticsteps.ManagerOption) *elasticsteps.Manager {
	return elasticsteps.NewManager(wrapClient(client), opts...)
}

// WithInstance adds a new es instance.
func WithInstance(name string, client *es7.Client) elasticsteps.ManagerOption {
	return elasticsteps.WithInstance(name, wrapClient(client))
}
