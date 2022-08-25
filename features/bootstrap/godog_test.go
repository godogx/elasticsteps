package bootstrap

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/godogx/expandvars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/godogx/elasticsteps"
)

// Used by init().
//
//nolint:gochecknoglobals
var (
	runGoDogTests bool

	out = new(bytes.Buffer)
	opt = godog.Options{
		Strict: true,
		Output: out,
	}
)

// This has to run on init to define -godog flag, otherwise "undefined flag" error happens.
//
//nolint:gochecknoinits
func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	godog.BindCommandLineFlags("", &opt)
}

func TestIntegration(t *testing.T) {
	if !runGoDogTests {
		t.Skip(`Missing "-godog" flag, skipping integration test.`)
	}

	flag.Parse()

	if opt.Randomize == 0 {
		opt.Randomize = rand.Int63() // nolint: gosec
	}

	drivers, err := newTestCases()
	require.NoError(t, err)

	for driver, manager := range drivers {
		driver, manager := driver, manager

		vars := expandvars.NewStepExpander(expandvars.Pairs{
			"DRIVER": driver,
		})

		t.Run(driver, func(t *testing.T) {
			RunSuite(t, "..", func(_ *testing.T, sc *godog.ScenarioContext) {
				vars.RegisterContext(sc)
				manager.RegisterContext(sc)
			})
		})
	}
}

func newTestCases() (map[string]*elasticsteps.Manager, error) {
	drivers := make(map[string]*elasticsteps.Manager)

	es7, err := newElasticsearch7(esAddr)
	if err != nil {
		return nil, fmt.Errorf("could not create a manager for %s: %w", typeES7, err)
	}

	drivers[typeES7] = es7

	return drivers, nil
}

func RunSuite(t *testing.T, path string, initScenario func(t *testing.T, ctx *godog.ScenarioContext)) {
	t.Helper()

	var paths []string

	files, err := os.ReadDir(filepath.Clean(path))
	assert.NoError(t, err)

	paths = make([]string, 0, len(files))

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".feature") {
			paths = append(paths, filepath.Join(path, f.Name()))
		}
	}

	for _, path := range paths {
		path := path

		t.Run(path, func(t *testing.T) {
			opt.Paths = []string{path}
			suite := godog.TestSuite{
				Name:                 "Integration",
				TestSuiteInitializer: nil,
				ScenarioInitializer: func(s *godog.ScenarioContext) {
					initScenario(t, s)
				},
				Options: &opt,
			}
			status := suite.Run()

			if status != 0 {
				fmt.Println(out.String())
				assert.Fail(t, "one or more scenarios failed in feature: "+path)
			}
		})
	}
}
