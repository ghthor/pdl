package database

import (
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	"log"
	"testing"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeAddFileExecutor)

	gospec.MainGoTest(r, t)
}
